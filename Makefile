include .envrc

# ----------------------------------------------
# searchly
# ----------------------------------------------

## build/searchly: build the cmd/searchly application
.PHONY: build/searchly
build/searchly:
	go build -ldflags="-s -w" -o=./bin/searchly ./cmd/searchly

# ----------------------------------------------
# deploy DBs
# ----------------------------------------------

## deploy/qdrant: create a qdrant docker container
.PHONY: deploy/qdrant
deploy/qdrant:
	-docker rm -f qdrant
	docker run --name qdrant --platform linux/arm64 -p 6333:6333 -d qdrant/qdrant:v1.12.1