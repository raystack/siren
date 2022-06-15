NAME="github.com/odpf/siren"
LAST_COMMIT := $(shell git rev-parse --short HEAD)
LAST_TAG := "$(shell git rev-list --tags --max-count=1)"
APP_VERSION := "$(shell git describe --tags ${LAST_TAG})-next"
PROTON_COMMIT := "ef83b9e9248e064a1c366da4fe07b3068266fe59"

.PHONY: all build test clean dist vet proto install

all: build

build: ## Build the siren binary
	@echo " > building siren version ${APP_VERSION}"
	go build -ldflags "-X main.Version=${APP_VERSION}" ${NAME}
	@echo " - build complete"

test: ## Run the tests
	go test ./... -coverprofile=coverage.out

coverage: ## Print code coverage
	go test -race -coverprofile coverage.out -covermode=atomic ./... && go tool cover -html=coverage.out

generate: ## run all go generate in the code base (including generating mock files)
	go generate ./...

lint: ## lint checker
	golangci-lint run

proto: ## Generate the protobuf files
	@echo " > generating protobuf from odpf/proton"
	@echo " > [info] make sure correct version of dependencies are installed using 'make install'"
	@buf generate https://github.com/odpf/proton/archive/${PROTON_COMMIT}.zip#strip_components=1 --template buf.gen.yaml --path odpf/siren
	@echo " > protobuf compilation finished"

clean: ## Clean the build artifacts
	rm -rf siren dist/

install: ## install required dependencies
	@echo "> installing dependencies"
	go mod tidy
	go install github.com/vektra/mockery/v2@v2.12.2
	go get -d google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go get github.com/golang/protobuf/proto@v1.5.2
	go get -d github.com/golang/protobuf/protoc-gen-go@v1.5.2
	go get google.golang.org/grpc@v1.40.0
	go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go get -d github.com/bufbuild/buf/cmd/buf@v0.54.1
	go get github.com/envoyproxy/protoc-gen-validate

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'