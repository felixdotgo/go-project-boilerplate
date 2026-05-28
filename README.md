# Golang project boilerplate
Starter template for Go monorepos/moduliths and microservices: Buf-managed protobufs and generated gRPC code, example services, a SolidJS auth frontend, and a local dev stack (Docker, Traefik, CoreDNS, mkcert).

## Project Structure

```
.
├── api/                # API specifications
│   ├── openapi/        # OpenAPI/Swagger definitions
│   └── proto/          # Protocol Buffer definitions
├── cmd/                # Application entry points
│   ├── foo/            # Example service
│   └── svc-auth/       # Authentication service
├── website/            # Frontend applications
│   └── auth/           # Authentication frontend (SolidStart/TypeScript)
├── pkg/                # Shared Go packages
│   ├── conv/           # Conversion utilities
│   ├── core/           # Core abstractions (server, repository, service base)
│   ├── log/            # Logging utilities
│   ├── migrator/       # Database migration tools
│   ├── rpc/            # Generated gRPC code (auto-generated from proto/)
│   └── tools/          # Build tools and dependencies
├── scripts/            # Build and setup scripts
├── devenv/             # Development environment configuration
└── init/               # Initialization files
```

### Key Directories

- **`cmd/`**: Contains the main applications. Each subdirectory represents a separate microservice
- **`pkg/`**: Reusable Go packages shared across services
- **`api/proto/`**: Protocol Buffer definitions for gRPC services
- **`pkg/rpc/`**: Auto-generated Go code from protobuf definitions
- **`website/`**: Web frontend applications (SolidStart/TypeScript)
- **`api/`**: API documentation and specifications

# Getting started
Before we start please make sure you have already installed these pieces of software
- Go >= 1.25.0
- Docker
- Cmake/Make
- mkcert

## Install
Run the following commands to start
```bash
# install necessary things
make install
# create and install self-signed certificate
make certs
# ramp up the environment
make docker.up
```

## Generate protobuf
You can define your protobuf inside `api/proto` directory and then run the following command to generate output
```bash
make generate-proto
```
All protobuf generated will be under `pkg/rpc` directory. To see how to implement API from generated code after run the command above, please refer to [Buf quick start](https://buf.build/docs/cli/quickstart/).

## Rename module path for forks
Forks can rewrite the template module path to match their repository.

```bash
# infer module path from git origin
make rename-module

# or provide it explicitly
make rename-module MODULE=github.com/your-org/your-repo
```

This maintenance utility lives under `scripts/`, not `cmd/`. It updates the Go module path, Go imports, protobuf `go_package` values, Buf config, and repository docs that still reference the previous module path.

## DNS setup
We use CoreDNS as a local DNS resolver

### For Mac
Set default DNS server to 127.0.0.1

### For Ubuntu
By default `systemd-resolve` will be used to resolve DNS and use port 53. So when we run `make run` the `core-dns` will failed.
To fix this issue, open `/etc/systemd/resolved.conf ` and add the following lines

```
DNSStubListener=no
DNS=127.0.0.1
```
then run
```
sudo systemctl restart systemd-resolved
```

### For Windows (WSL2)
Create `.wslconfig` with the following content
```
[wsl2]
networkingMode=default
localhostForwarding=true
```
Shutdown and restart WSL2

In Docker Desktop, go to `Settings > Resources > WSL Integration`, then disable and re-enable integration for your WSL2 distro.

Uncomment the `insecure: true` line in `devenv/traefik/traefik.yml` to enable insecure access to the Traefik dashboard.
