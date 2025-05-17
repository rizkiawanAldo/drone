.PHONY: init generate build test test-unit run docker-build docker-run clean

# Variables
GOPATH := $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
OAPI_CODEGEN := $(BIN_DIR)/oapi-codegen
GENERATED_DIR := ./generated
SERVER_DIR := ./cmd/server
SERVER_BINARY := server

init:
	go mod tidy
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
	docker-compose build

docker-run:
	docker-compose up

clean:
	rm -f $(SERVER_BINARY)
	rm -f coverage.out 