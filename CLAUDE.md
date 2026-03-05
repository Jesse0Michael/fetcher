# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
make run                    # go run ./cmd/fetcher

# Build
make build                  # produces bin/fetcher (HTTP server)

# Test and lint
make test                   # go test -cover ./internal/... && golangci-lint run ./internal/...
make test-coverage          # generates HTML coverage report

# Run a single test
go test ./internal/service/... -run TestFunctionName
go test ./internal/server/... -run TestFunctionName

# Docker
make build-docker
make run-docker             # runs on port 8080

# Generate OpenAPI client + go generate
make gen
```

Configuration is loaded from a `.env` file (if present) and environment variables via `envconfig`. Log level is controlled by `LOG_LEVEL` env var (DEBUG/INFO/WARN/ERROR).

## Architecture

**Fetcher** is a Go HTTP API that aggregates social media feeds from multiple platforms into a normalized `FeedItem` format, sorted by timestamp.

### Request Flow

1. `GET /feed?twitterID=...&bloggerID=...&...` hits `server.feed()` handler
2. Handler decodes query params into `service.FetcherRequest` using `go-request`
3. `service.Fetcher.Feeds()` fans out concurrent goroutines — one per requested platform ID
4. Results are aggregated, sorted descending by `TS` (Unix timestamp), and returned as `{"items": [...]}`

There is also a `GET /proxy?url=...` endpoint that proxies arbitrary HTTP GET requests (used to work around CORS for some feed sources).

### Key Packages

- **`internal/service/`** — core business logic. Each platform has its own file (`twitter.go`, `blogger.go`, `swarm.go`, `soundcloud.go`, `deviantart.go`, `untappd.go`, `instagram.go`). Each implements the `Feeder` interface: `Feed(ctx, id) ([]FeedItem, error)`.
- **`internal/server/`** — HTTP layer using `gorilla/mux`. `server.go` sets up routes; `feed.go` and `proxy.go` are the two handlers.
- **`cmd/fetcher/main.go`** — entry point; wires `service.Fetcher` + `server.Server` with config from env.

### Adding a New Platform

1. Create `internal/service/<platform>.go` implementing the `Feeder` interface
2. Add a config struct with `envconfig` tags and include it in `service.Config`
3. Instantiate it in `NewFetcher()` and add a `req.<PlatformID> != ""` branch in `Feeds()`
4. Add a query param field to `FetcherRequest`

### Notes

- Instagram is currently disabled (commented out in `NewFetcher`).
- The Makefile references a `cmd/lambda/` build target, but that directory does not exist.
- `FeedItem.Content` may contain HTML (e.g., Twitter handler generates HTML markup).
- API spec lives in `api/openapi.yaml`; docs are generated with `make gen-docs`.
