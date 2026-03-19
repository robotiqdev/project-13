#!/bin/bash
set -euo pipefail

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo 'none')
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-X github.com/app/service/internal/version.Version=${VERSION} -X github.com/app/service/internal/version.Commit=${COMMIT} -X github.com/app/service/internal/version.BuildDate=${BUILD_DATE}"

go build -ldflags "${LDFLAGS}" -o bin/server ./cmd/server/
