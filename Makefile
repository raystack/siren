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
	go list ./... | grep -v extern | xargs go test -count 1 -cover -race -timeout 1m

dist:
	@bash ./scripts/build.sh
