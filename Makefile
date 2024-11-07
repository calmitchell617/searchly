include .envrc

# ----------------------------------------------
# BUILD APPS
# ----------------------------------------------

## build/benchmark: build the cmd/benchmark application
.PHONY: build/benchmark
build/benchmark:
	go build -ldflags="-s -w" -o=./bin/benchmark ./cmd/benchmark

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
	./bin/prepare -system milvus -host ${MILVUS_HOST} -num-entities 1000000

## prepare/weaviate
.PHONY: prepare/weaviate
prepare/weaviate: build/prepare
	./bin/prepare -system weaviate -host ${WEAVIATE_HOST} -num-entities 1000000

## prepare/qdrant
.PHONY: prepare/qdrant
prepare/qdrant: build/prepare
	./bin/prepare -system qdrant -host ${QDRANT_HOST} -num-entities 1000000

# ----------------------------------------------
# RUN BENCHMARKS
# ----------------------------------------------

## benchmark/milvus
.PHONY: benchmark/milvus
benchmark/milvus: build/benchmark
	./bin/benchmark -system milvus -host ${MILVUS_HOST} -duration 61m

## benchmark/weaviate
.PHONY: benchmark/weaviate
benchmark/weaviate: build/benchmark
	./bin/benchmark -system weaviate -host ${WEAVIATE_HOST} -duration 61m

## benchmark/qdrant
.PHONY: benchmark/qdrant
benchmark/qdrant: build/benchmark
	./bin/benchmark -system qdrant -host ${QDRANT_HOST} -duration 61m

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