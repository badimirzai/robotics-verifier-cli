# robostack-cli Makefile
BINARY := robostack
PKG := ./...
GOFLAGS :=

.PHONY: help tidy fmt vet test lint build run validate clean

help:
	@echo "Targets:"
	@echo "  tidy      - go mod tidy"
	@echo "  fmt       - gofmt all go files"
	@echo "  vet       - go vet"
	@echo "  test      - run unit tests"
	@echo "  lint      - golangci-lint (if installed)"
	@echo "  build     - build ./cmd binary into ./bin"
	@echo "  run       - run CLI (pass ARGS=\"...\")"
	@echo "  validate  - run validate on example spec"
	@echo "  clean     - remove ./bin"

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet $(PKG)

test:
	go test -race $(PKG)

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || \
		(echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; exit 1)

build:
	mkdir -p bin
	go build $(GOFLAGS) -o bin/$(BINARY) .

run:
	go run . $(ARGS)

validate:
	go run . validate -f examples/amr_basic.yaml

clean:
	rm -rf bin