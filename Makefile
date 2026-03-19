.PHONY: build run test lint

build:
	scripts/build.sh \
		-ldflags "-X github.com/robotiqdev/project-13/internal/version.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev) \
		          -X github.com/robotiqdev/project-13/internal/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown) \
		          -X github.com/robotiqdev/project-13/internal/version.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" \
		-o bin/server ./cmd/server/

run:
	go run \
		-ldflags "-X github.com/robotiqdev/project-13/internal/version.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev) \
		          -X github.com/robotiqdev/project-13/internal/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown) \
		          -X github.com/robotiqdev/project-13/internal/version.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" \
		./cmd/server/

test:
	go test ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping lint"; \
	fi
