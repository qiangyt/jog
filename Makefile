GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
SERVER_PROTO_FILES=$(shell find server -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/go-kratos/kratos/cmd/kratos/v2@v2.0.0-20220128070526-34d0cccefd7b
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@v2.0.0-20220128070526-34d0cccefd7b
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@v2.0.0-20220128070526-34d0cccefd7b
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.1

.PHONY: errors
# generate errors code
errors:
	protoc --proto_path=. \
               --proto_path=./third_party \
               --go_out=paths=source_relative:. \
               --go-errors_out=paths=source_relative:. \
               $(API_PROTO_FILES)

.PHONY: config
# generate server proto
config:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
	       $(SERVER_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
 	       --openapi_out==paths=source_relative:. \
	       $(API_PROTO_FILES)

.PHONY: clean
# clean
clean:
	rm -rf ./target

.PHONY: build
# build
build:
	mkdir -p target/
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/linux/jog ./cmd
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/darwin/jog.amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/darwin/jog.arm64 ./cmd
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/windows/jog.exe ./cmd

.PHONY: generate
# generate
generate:
	statik -src=./web/dist -f
	go generate ./...
	go fmt ./...

.PHONY: test
# test
test:
	# go test ./... -v -covermode=count -coverprofile=coverage.out gcflags=all=-l
	go test ./... -covermode=count -coverprofile=coverage.out gcflags=all=-l
	go tool cover -html=./coverage.out -o ./coverage.html
	#go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest
	#https://github.com/alecthomas/gometalinter/releases/tag/v3.0.0
	#goreportcard-cli -v
	goreportcard-cli

.PHONY: all
# generate all
all:
	make api;
	make errors;
	make config;
	make generate;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
