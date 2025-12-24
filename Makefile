# robostack-cli Makefile
BINARY := robostack
PKG := ./...
GOFLAGS :=

# Optional overrides:
#   make run ARGS="validate -f examples/amr_basic.yaml"
#   make validate FILE=examples/amr_basic.yaml
ARGS ?=
FILE ?= examples/amr_basic.yaml

.PHONY: help tidy fmt vet test lint build install run validate clean

help:
	@echo "Targets:"
	@echo "  tidy       - go mod tidy"
	@echo "  fmt        - gofmt all go files"
	@echo "  vet        - go vet"
	@echo "  test       - run unit tests"
	@echo "  lint       - golangci-lint (if installed)"
	@echo "  build      - build binary into ./bin/$(BINARY)"
	@echo "  install    - install binary into $(shell go env GOPATH)/bin"
	@echo "  run        - run CLI (requires ARGS="...")"
	@echo "  validate   - run validate on FILE (default: $(FILE))"
	@echo "  clean      - remove ./bin"
	@echo ""
	@echo "Examples:"
	@echo "  make validate"
	@echo "  make validate FILE=examples/amr_basic.yaml"
	@echo "  make run ARGS=\"validate -f examples/amr_basic.yaml\""
	@echo "  make build && ./bin/$(BINARY) validate -f examples/amr_basic.yaml"

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet $(PKG)

test:
	go test $(PKG)

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || 		(echo "golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; exit 1)

build:
	mkdir -p bin
	go build $(GOFLAGS) -o bin/$(BINARY) .

install:
	go install $(GOFLAGS) .
	@echo ""
	@echo "Installed to: $$(go env GOPATH)/bin/$(BINARY)"
	@echo "If '$(BINARY)' is not found, add this to your PATH:"
	@echo "  export PATH=\"$$(go env GOPATH)/bin:$$PATH\""

run:
	@if [ -z "$(strip $(ARGS))" ]; then 		echo "ERROR: ARGS is required."; 		echo "Example: make run ARGS=\"validate -f examples/amr_basic.yaml\""; 		exit 2; 	fi
	go run . $(ARGS)

validate:
	go run . validate -f $(FILE)

clean:
	rm -rf bin
