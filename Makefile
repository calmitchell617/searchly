include .envrc

# ----------------------------------------------
# BUILD APPS
# ----------------------------------------------

## build/searchly: build the cmd/searchly application
.PHONY: build/searchly
build/searchly:
	go build -ldflags="-s -w" -o=./bin/searchly ./cmd/searchly

## build/prepare: build the cmd/prepare application
.PHONY: build/prepare
build/prepare:
	go build -ldflags="-s -w" -o=./bin/prepare ./cmd/prepare

# ----------------------------------------------
# PREPARE DBS
# ----------------------------------------------

## prepare/milvus
.PHONY: prepare/milvus
prepare/milvus: build/prepare
	./bin/prepare -system milvus -num-entities 1000000 -host ${MILVUS_HOST}

## prepare/qdrant
.PHONY: prepare/qdrant
prepare/qdrant: build/prepare
	./bin/prepare -system qdrant

## prepare/weaviate
.PHONY: prepare/weaviate
prepare/weaviate: build/prepare
	./bin/prepare -system weaviate

# ----------------------------------------------
# DEPLOY DBS
# ----------------------------------------------

## deploy/qdrant: create a qdrant docker container
.PHONY: deploy/qdrant
deploy/qdrant:
	-docker rm -f qdrant
	docker run --name qdrant --platform linux/amd64 -p 6333:6333 -p 6334:6334 -d qdrant/qdrant:v1.12.1
	docker logs -f qdrant

## deploy/weaviate: create a weaviate docker container
.PHONY: deploy/weaviate
deploy/weaviate:
	-docker rm -f weaviate
	docker run --name weaviate -p 8080:8080 -p 50051:50051 --platform linux/amd64 -d cr.weaviate.io/semitechnologies/weaviate:1.27.1
	docker logs -f weaviate

## deploy/milvus: create a milvus docker container using the docker compose file ./setup/milvus/docker-compose.yml
.PHONY: deploy/milvus
deploy/milvus:
	-docker compose -v -f setup/milvus/docker-compose.yml down
	docker compose -f setup/milvus/docker-compose.yml up -d
	docker compose -f setup/milvus/docker-compose.yml logs -f

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor