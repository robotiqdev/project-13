# Service

## Prerequisites

- Go 1.21+

## Building

```bash
make build
```

Or manually with ldflags:

```bash
go build -ldflags "-X github.com/app/service/internal/version.Version=v1.2.3 \
  -X github.com/app/service/internal/version.Commit=abc1234 \
  -X github.com/app/service/internal/version.BuildDate=2024-01-15T10:30:00Z" \
  -o bin/server ./cmd/server/
```

## Running

```bash
make run
```

Or directly:

```bash
./bin/server
```

## Testing

```bash
make test
```

Or:

```bash
go test ./...
```

## Endpoints

- `GET /health` — returns service health status
- `GET /version` — returns build metadata

## Version Injection

Build metadata is injected at compile time using Go's `-ldflags` mechanism. The `scripts/build.sh` script sets the following variables:

- `version.Version` — from `git describe --tags`
- `version.Commit` — from `git rev-parse --short HEAD`
- `version.BuildDate` — current UTC time

Example response from `/version`:

```bash
curl localhost:8080/version | jq .
```

```json
{
  "version": "v1.2.3",
  "commit": "abc1234",
  "buildDate": "2024-01-15T10:30:00Z"
}
```

When built without ldflags, all fields default to empty strings (`""`).
