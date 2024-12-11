PROJECT_NAME=grpc-todo
MODULE_NAME=grpc-todo

.DEFAULT_GOAL := build

.PHONY: proto
proto:
	@protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=pkg/api/v1 --go_opt=paths=source_relative --go-grpc_out=pkg/api/v1 --go-grpc_opt=paths=source_relative api/proto/v1/*.proto
	@protoc --proto_path=api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api api/proto/v1/*.proto
	@protoc --proto_path=api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 api/proto/v1/*.proto
.PHONY: serve
serve:
	@go run cmd/server/main.go

.PHONY: client_rpc
client_rpc:
	@go run cmd/client-grpc/main.go

.PHONY: client_rest
client_rest:
	@go run cmd/client-rest/main.go

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