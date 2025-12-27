# architon-cli Makefile
BINARY := architon-cli
PKG := ./...
GOFLAGS :=

# Optional overrides:
#   make run ARGS="check examples/amr_basic.yaml"
#   make check FILE=examples/amr_basic.yaml
ARGS ?=
FILE ?= examples/amr_parts.yaml

.PHONY: help tidy fmt vet test lint build install run check validate clean

help:
	@echo "Targets:"
	@echo "  tidy       - go mod tidy"
	@echo "  fmt        - gofmt all go files"
	@echo "  vet        - go vet"
	@echo "  test       - run unit tests"
	@echo "  lint       - golangci-lint (if installed)"
	@echo "  build      - build binary into ./bin/$(BINARY)"
	@echo "  install    - install binary into $(shell go env GOPATH)/bin (and arch symlink)"
	@echo "  run        - run CLI (requires ARGS="...")"
	@echo "  check      - run check on FILE (default: $(FILE))"
	@echo "  validate   - alias for check"
	@echo "  clean      - remove ./bin"
	@echo ""
	@echo "Examples:"
	@echo "  make check"
	@echo "  make check FILE=examples/amr_basic.yaml"
	@echo "  make run ARGS=\"check examples/amr_basic.yaml\""
	@echo "  make build && ./bin/$(BINARY) check examples/amr_basic.yaml"

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
	ln -sf "$$(go env GOPATH)/bin/$(BINARY)" "$$(go env GOPATH)/bin/arch"
	@echo ""
	@echo "Installed to: $$(go env GOPATH)/bin/$(BINARY)"
	@echo "Symlinked: $$(go env GOPATH)/bin/arch"
	@echo "If 'arch' is not found, add this to your PATH:"
	@echo "  export PATH=\"$$(go env GOPATH)/bin:$$PATH\""

run:
	@if [ -z "$(strip $(ARGS))" ]; then 		echo "ERROR: ARGS is required."; 		echo "Example: make run ARGS=\"check examples/amr_basic.yaml\""; 		exit 2; 	fi
	go run . $(ARGS)

check:
	go run . check $(FILE)

validate:
	go run . check $(FILE)

clean:
	rm -rf bin
