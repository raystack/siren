NAME="github.com/goto/siren"
LAST_COMMIT := $(shell git rev-parse --short HEAD)
LAST_TAG := "$(shell git rev-list --tags --max-count=1)"
APP_VERSION := "$(shell git describe --tags ${LAST_TAG})-next"
PROTON_COMMIT := "c4c8fa81f81e7f76ba4d0569b09b016234e01915"

.PHONY: all build test clean dist vet proto install

all: build


build: build-main build-plugins

build-main: ## Build the siren binary
	@echo " > building siren version ${APP_VERSION}"
	go build -ldflags "-X main.Version=${APP_VERSION}" ${NAME}
	@echo " > building plugins version ${APP_VERSION}"
	@echo " - build complete"

build-plugins:
	@echo " > building plugins"
	@go list -m | grep providers | while read path; do go build -o ./plugin/${shell basename "$$path"} "$$path"; done
	@echo " - build complete"

test: ## Run the tests
	go test -race $(shell go list ./... | grep -v /test/) -covermode=atomic -coverprofile=coverage.out

e2e-test: build-plugins ## Run all e2e tests
	go test -v -race ./test/e2e_test/... -coverprofile=coverage.out --timeout 300s

coverage: ## Print code coverage
	go test -race -coverprofile coverage.out -covermode=atomic ./... && go tool cover -html=coverage.out

generate: ## run all go generate in the code base (including generating mock files)
	@go generate ./...
	@echo " > generating mock files"
	@mockery

lint: ## lint checker
	golangci-lint run

proto: ## Generate the protobuf files
	@echo " > generating protobuf from gotocompany/proton"
	@echo " > [info] make sure correct version of dependencies are installed using 'make install'"
	@buf generate https://github.com/goto/proton/archive/${PROTON_COMMIT}.zip#strip_components=1 --template buf.gen.yaml --path gotocompany/siren
	@echo " > protobuf compilation finished"

clean: ## Clean the build artifacts
	rm -rf siren dist/

update-swagger-md:
	@echo "> updating reference api docs"
	@npx swagger-markdown -i proto/siren.swagger.yaml -o docs/docs/reference/api.md

install: ## install required dependencies
	go get -d github.com/vektra/mockery/v2@v2.40.1
	go get -d google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	go get google.golang.org/protobuf/proto@v1.32.0
	go get google.golang.org/grpc@v1.61.0
	go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	go get -d github.com/bufbuild/buf/cmd/buf@v1.29.0
	go get github.com/envoyproxy/protoc-gen-validate@v1.0.4

clean-doc:
	@echo "> cleaning up auto-generated docs"
	@rm -rf ./docs/docs/reference/cli
	@rm -f ./docs/docs/reference/api.md

# Generates the config file documentation.
# remove ansi color & escape html
doc: clean-doc update-swagger-md
	@echo "> generate cli docs"
	@go run . reference --plain | sed '1 s,.*,# CLI,' > ./docs/docs/reference/cli.md

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
