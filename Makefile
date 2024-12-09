PROJECT_NAME=grpc-todo
MODULE_NAME=grpc-todo

.DEFAULT_GOAL := build

.PHONY: proto
proto:
	@protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=api/proto/v1 --go_opt=paths=source_relative --go-grpc_out=api/proto/v1 --go-grpc_opt=paths=source_relative api/proto/v1/*.proto
.PHONY: run
run:
	@go run cmd/grpc-clean/main.go

.PHONY: build
build:
	@go build ./cmd/grpc-clean/

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test -v -coverprofile coverage.out ./...

.PHONY: coverage
coverage:
	@go tool cover -html=coverage.out

.PHONY: get
get:
	@go mod download