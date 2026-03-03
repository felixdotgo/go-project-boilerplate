# Copilot Instructions for go-project-boilerplate

## Repository Overview

This is a Go monorepo boilerplate for building microservices. It provides:
- Buf-managed Protocol Buffer definitions and generated gRPC/HTTP gateway code
- An example authentication service (`svc-auth`) with JWT, refresh tokens, and OAuth2 (Google, Facebook, GitHub)
- A local development stack (Docker, Traefik, CoreDNS, PostgreSQL)
- A SolidJS/TypeScript authentication frontend (`website/auth`)
- Shared Go packages in `pkg/`

**Module path:** `github.com/0x46656C6978/go-project-boilerplate`

## Repository Structure

```
.
├── api/
│   ├── openapi/        # Generated OpenAPI/Swagger specs
│   └── proto/          # Protocol Buffer source definitions
│       ├── api/        # Service API protos
│       ├── common/     # Shared proto types
│       └── models/     # Proto model types
├── cmd/
│   ├── foo/            # Minimal example service
│   └── svc-auth/       # Full authentication service
│       ├── cmd/migrator/  # DB migration runner
│       ├── config/        # Viper-based config (reads .env)
│       ├── entity/        # Domain entities (pure Go structs)
│       ├── httpapi/       # gRPC-Gateway HTTP handlers
│       ├── migrations/sql/# SQL migration files
│       ├── oauth/         # OAuth2 provider implementations
│       ├── repository/    # Data access layer (GORM + PostgreSQL)
│       └── service/       # Business logic layer
├── devenv/             # Docker Compose dev environment
│   ├── coredns/        # Local DNS config (*.goproject.localhost)
│   └── traefik/        # Reverse proxy config
├── pkg/
│   ├── conv/           # Type conversion utilities
│   ├── core/           # Base structs: HandlerServer, ServiceBase, RepositoryBase
│   ├── log/            # Zerolog-based structured logger
│   ├── migrator/       # Database migration runner
│   ├── rpc/            # Auto-generated gRPC/gateway code (DO NOT edit manually)
│   └── tools/          # Build tool dependencies (tools.go pattern)
├── scripts/            # setup_environment.sh and helpers
└── website/auth/       # SolidStart/TypeScript auth frontend
```

## Build, Run, and Test

### Prerequisites
- Go >= 1.25 (exact version in `.go-version`: `go1.25`)
- Docker + Docker Compose
- `make` / `cmake`
- `mkcert` (for local TLS certificates)
- `buf` CLI (for protobuf generation)
- `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway` (installed via `make install`)

### Common Make Commands

```bash
make install          # Install all dev dependencies (brew, protoc, buf, Go tools)
make certs            # Generate self-signed TLS certs for *.goproject.localhost
make docker.up        # Start dev environment (PostgreSQL, Traefik, CoreDNS)
make docker.down      # Stop and clean up dev containers

make run CMD=svc-auth # Run a service locally (go run cmd/<CMD>/main.go)
make build CMD=svc-auth  # Build Docker image for a service
make up CMD=svc-auth     # Build + docker compose up (attached)
make upd CMD=svc-auth    # Build + docker compose up -d (detached)

make generate-proto   # Regenerate code from .proto files (outputs to pkg/rpc/ and api/openapi/)
make vendor           # Run go mod vendor
```

### Running Tests

There are currently no automated test files in the repository. If adding tests:
- Use Go's standard `testing` package
- Place test files next to the code they test with `_test.go` suffix
- Run with `go test ./...` from the repo root

### Building

```bash
# Build a specific service binary
go build ./cmd/svc-auth/...

# Build all packages
go build ./...
```

## Architecture Patterns

### Adding a New Service

1. Create `cmd/<service-name>/` with:
   - `main.go` — wire up config, logger, DB, `core.HandlerServer`, register gRPC handlers
   - `config/config.go` — Viper config struct reading from `.env`
   - `Dockerfile` + `docker-compose.yaml`
2. Define protobuf in `api/proto/api/<service>/v1/`
3. Run `make generate-proto` to produce Go code in `pkg/rpc/`
4. Implement the generated `<Service>Server` interface in `httpapi/`
5. Register the handler via `<service>v1.Register<Service>HandlerServer(ctx, server.GetMux(), impl)`

### Layered Architecture (svc-auth pattern)

```
httpapi/   → implements generated gRPC server interface, calls service layer
service/   → business logic, depends on repository interface
repository/→ data access, implements repository interface using GORM
entity/    → domain structs (no framework dependencies)
model/     → GORM model structs (database representation)
```

- **Interfaces first**: define `UserServiceInterface` and `UserRepoInterface` in their respective packages before implementing
- **entity/** = domain layer (pure Go, no ORM tags)
- **repository/model/** = persistence layer (GORM models with `gorm:` tags)
- **repository/dto/** = data transfer objects for repo layer

### Core Abstractions (`pkg/core`)

- `core.HandlerServer` — wraps `grpc-gateway` `runtime.ServeMux` with HTTP/2 + h2c support; use `server.EnableCORS()` for browser clients
- `core.ServiceBase` — embeds a named zerolog logger; embed this in every service struct
- `core.RepositoryBase` — embeds a named zerolog logger; embed this in every repository struct

### Configuration

- Config is loaded from a `.env` file via Viper (`AutomaticEnv` + `ReadInConfig`)
- Each service has its own `config/config.go` with a `Config` struct using `mapstructure` tags matching env var names (all lowercase with underscores)
- See `cmd/svc-auth/.env.example` for the full list of variables

### Protobuf / gRPC-Gateway

- Proto files live in `api/proto/`
- Generated code goes to `pkg/rpc/` — **never edit files there manually**
- HTTP bindings are defined inline in `.proto` files using `google.api.http` options
- Regenerate with `make generate-proto` (requires buf and local protoc plugins)
- Buf config: `buf.yaml` (modules + lint rules) and `buf.gen.yaml` (code generation plugins)

### Database / Migrations

- PostgreSQL via GORM (`gorm.io/driver/postgres`)
- Migrations stored as SQL files in `cmd/<service>/migrations/sql/`
- Migration runner: `cmd/svc-auth/cmd/migrator/migrator.go` uses `pkg/migrator`
- Run migrations: `go run cmd/svc-auth/cmd/migrator/migrator.go`
- Dev DB: `postgres:18` container (user: `postgres`, password: `root`, db: `postgres`, port: `5432`)

### OAuth2

- Supported providers: Google, Facebook, GitHub
- Each provider implements the `oauth.Provider` interface in `cmd/svc-auth/oauth/`
- OAuth tokens stored encrypted (AES-256-GCM) in the database — requires a 32-byte `OAUTH_ENCRYPTION_KEY`
- `oauth_encryption_key` in config must be exactly 32 bytes

### Logging

- Uses `pkg/log` (zerolog wrapper)
- `log.NewLogger(debug bool)` — pass `true` for debug output
- Log with context: `logger.With("key", value).Info("message")`
- Services and repositories get a named logger via `ServiceBase`/`RepositoryBase`

## Code Style

- Follow `.editorconfig`: tabs for Go/Makefile, spaces (2) for YAML/JSON/other
- Go indentation: tabs (enforced by `gofmt`)
- Interfaces are named `<Type>Interface` (e.g., `UserServiceInterface`, `UserRepoInterface`)
- Constructor functions: `New<Type>(...)` returning the interface type
- Nil-guard all receiver methods (check `if x == nil { return zero }`)
- Error constants defined in `service/common.go` using `errors.New(...)`

## Known Issues and Workarounds

### DNS (CoreDNS on port 53)

On Ubuntu, `systemd-resolved` occupies port 53. Fix:
```
# /etc/systemd/resolved.conf
DNSStubListener=no
DNS=127.0.0.1
```
Then `sudo systemctl restart systemd-resolved`.

On macOS: set DNS server to `127.0.0.1`.

On WSL2: configure `.wslconfig` with `networkingMode=default` + `localhostForwarding=true`, restart WSL2, re-enable Docker Desktop WSL integration.

### Traefik Dashboard (WSL2)

Uncomment `insecure: true` in `devenv/traefik/traefik.yml` to enable the dashboard over HTTP.

### Protobuf Generation Requires Local Tools

`make generate-proto` requires `protoc-gen-go`, `protoc-gen-go-grpc`, and `protoc-gen-grpc-gateway` to be in `$PATH`. Run `make install` first.

## CI / Workflow Notes

There are no GitHub Actions workflows currently configured. When adding CI:
- Use `go build ./...` to verify compilation
- Use `go vet ./...` for static analysis
- Use `go test ./...` for tests
- Docker images are built per-service with `make build CMD=<service>`
