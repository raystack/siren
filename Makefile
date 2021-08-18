NAME="github.com/odpf/siren"
VERSION=$(shell git describe --always --tags 2>/dev/null)
COVERFILE="/tmp/siren.coverprofile"

.PHONY: all build clean

all: build

build:
	go build -ldflags "-X main.Version=${VERSION}" ${NAME}

clean:
	rm -rf siren dist/

test:
	go test ./... -coverprofile=coverage.out

test-coverage: test
	go tool cover -html=coverage.out

dist:
	@bash ./scripts/build.sh

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on go mod vendor  && swagger generate spec -o ./api/handlers/swagger.yaml --scan-models

swagger-serve: check-swagger
	swagger serve -F=swagger api/handlers/swagger.yaml

generate-proto: ## regenerate protos
	@echo " > cloning protobuf from odpf/proton"
	@echo " > generating protobuf"
	@buf generate --template buf.gen.yaml https://github.com/odpf/proton.git --path odpf/odin
	@echo " > protobuf compilation finished"