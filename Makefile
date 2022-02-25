GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
SERVER_PROTO_FILES=$(shell find ./pkg/server/conf -name *.proto)
#API_PROTO_FILES=$(shell find ./api/proto -name *.proto)
API_PROTO_FILES=./api/proto/error_reason.proto ./api/proto/index.proto ./api/proto/greeter.proto

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/go-kratos/kratos/cmd/kratos/v2@v2.0.0-20220128070526-34d0cccefd7b
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@v2.0.0-20220128070526-34d0cccefd7b
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@v2.0.0-20220128070526-34d0cccefd7b
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.1
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/rakyll/statik@v0.1.7

	cd ./web && npm install


.PHONY: web-dist
# web-dist
web-dist:
	cd ./web && quasar build
	statik -src=./web/dist/spa -dest=./web -f -include=* -ns=web

.PHONY: embed-res
# embed-res
embed-res:
	statik -src=./res/raw -dest=./res -f -include=* -ns=res

.PHONY: protoc
# protoc
protoc:
	protoc --proto_path=. \
	       --proto_path=./api/proto \
	       --proto_path=./third_party \
 	       --go_out=./api/go \
 	       --go-http_out=./api/go \
 	       --go-grpc_out=./api/go \
 	       --openapi_out=./api \
	       $(API_PROTO_FILES)

	protoc --proto_path=. \
	       --proto_path=./api/proto \
         --proto_path=./third_party \
         --go_out=./api/go \
         --go-errors_out=./api/go \
         $(API_PROTO_FILES)

	protoc --proto_path=. \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:. \
	       $(SERVER_PROTO_FILES)

.PHONY: generate
# generate server & errors code & api
generate:
	make protoc;
	make embed-res;

	go generate ./...
	go fmt ./...

.PHONY: clean
# clean
clean:
	rm -rf ./target

.PHONY: linux-amd64
# build for linux/amd64
linux-amd64:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog.linux-amd64 ./

.PHONY: linux-arm64
# build for linux/arm64
linux-arm64:
	GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog.linux-arm64 ./

.PHONY: darwin-amd64
# build for darwin/amd64
darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog.darwin-amd64 ./

.PHONY: darwin-arm64
# build for darwin/arm64
darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog.darwin-arm64 ./

.PHONY: windows-amd64
# build for windows/amd64
windows-amd64:
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog.amd64.exe ./

.PHONY: windows-amd64
# build for windows/arm64
windows-arm64:
	GOOS=windows GOARCH=arm64 go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog.arm64.exe ./

.PHONY: build
# build
build:
	go build -trimpath -ldflags "-X main.Version=$(VERSION)" -o ./target/jog ./

.PHONY: release
# release
release:
	make linux-amd64;
	make linux-arm64;
	make darwin-amd64;
	make darwin-arm64;
	make windows-amd64;
	make windows-arm64;

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
# all
all:
	make clean;
	make web-dist;
	make generate;
	make build;

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
