.PHONY: init generate build test test-unit run docker-build docker-run docker-build-single docker-run-single clean

# Variables
GOPATH := $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
OAPI_CODEGEN := $(BIN_DIR)/oapi-codegen
GENERATED_DIR := ./generated
SERVER_DIR := ./cmd/server
SERVER_BINARY := server
DOCKER_IMAGE := drone-app

init:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install github.com/golang/mock/mockgen@latest
	mkdir -p $(GENERATED_DIR)
	$(MAKE) generate

generate:
	$(OAPI_CODEGEN) -package generated -generate types,server,spec api.yaml > $(GENERATED_DIR)/api.gen.go

build:
	go build -o $(SERVER_BINARY) $(SERVER_DIR)/main.go

test: test-unit

test-unit:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

run: build
	./$(SERVER_BINARY)

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

clean:
	rm -f $(SERVER_BINARY)
	rm -f coverage.out 

reset:
	rm -f $(SERVER_BINARY)
	rm -f coverage.out
	rm -rf $(GENERATED_DIR)