VERSION=`git describe --tags --abbrev=0`
COMPILED=`date -u +%Y%m%d-%H%M%S`
LDFLAGS="-s -w -X github.com/riussi/4sq-exports/cmd.compiled=$(COMPILED) -X github.com/riussi/4sq-exports/cmd.version=$(VERSION)"
GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: build-osx

build-all: build-osx build-linux build-windows

build-osx: $(GOFILES)
	rm -rf 4sq-exports
	CGOENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS)
#	mv 4sq-exports 4sq-exports-osx-$(VERSIONSTRING)

build-linux: $(GOFILES)
	CGOENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS)
	mv 4sq-exports 4sq-exports-linux-$(VERSIONSTRING)

build-windows: $(GOFILES)
	CGOENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS)
	mv 4sq-exports.exe 4sq-exports-win-$(VERSIONSTRING).exe

test: test-all

test-all:
	CGOENABLED=0 go vet $(GOPACKAGES)
	CGOENABLED=0 go test $(GOPACKAGES)

lint: lint-all

lint-all:
	CGOENABLED=0 go fmt $(GOPACKAGES)
	CGOENABLED=0 golint $(GOPACKAGES)
	CGOENABLED=0 gometalinter $(GOPACKAGES)
	@golint -set_exit_status $(GOPACKAGES)
