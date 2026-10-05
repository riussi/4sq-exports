VERSION=`git describe --tags --abbrev=0`
COMPILED=`date -u +%Y%m%d-%H%M%S`
LDFLAGS="-s -w -X github.com/riussi/4sq-exports/cmd.clientID=$(FOURSQCLIENTID) -X github.com/riussi/4sq-exports/cmd.clientSecret=$(FOURSQCLIENTSECRET) -X github.com/riussi/4sq-exports/cmd.compiled=$(COMPILED) -X github.com/riussi/4sq-exports/cmd.version=$(VERSION)"
GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')

default: build-osx

build-all: build-osx build-linux build-windows

build-osx: $(GOFILES)
	rm -rf 4sq-exports
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS)
#	mv 4sq-exports 4sq-exports-osx-$(VERSION)

build-linux: $(GOFILES)
	rm -rf 4sq-exports
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS)
	mv 4sq-exports 4sq-exports-linux-$(VERSION)

build-windows: $(GOFILES)
	rm -rf 4sq-exports
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS)
	mv 4sq-exports.exe 4sq-exports-win-$(VERSION).exe

test: test-all

test-all:
	go vet ./...
	go test ./...

lint: lint-all

lint-all:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && echo "gofmt: files need formatting" && exit 1)
	go vet ./...
	go tool staticcheck ./...

vuln:
	go tool govulncheck ./...

.PHONY: default build-all build-osx build-linux build-windows test test-all lint lint-all vuln
