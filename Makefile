.PHONY: build

MAINFILE := cmd/tui/main.go
BINNAME := algorand-navigator

GOLDFLAGS += -X github.com/winder/algorand-navigator/version.Hash=$(shell git log -n 1 --pretty="%H")
GOLDFLAGS += -X github.com/winder/algorand-navigator/version.ShortHash=$(shell git log -n 1 --pretty="%h")
GOLDFLAGS += -X github.com/winder/algorand-navigator/version.CompileTime=$(shell date -u +%Y-%m-%dT%H:%M:%S%z)
GOLDFLAGS += -X "github.com/winder/algorand-navigator/version.ReleaseVersion=Dev Build"

build:
	go build -o $(BINNAME) -ldflags='${GOLDFLAGS}' $(MAINFILE)

fmt:
	go fmt ./...

lint:
	golint -set_exit_status ./...
	go vet ./...
	golangci-lint run

dep:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
