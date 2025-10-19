# Golang project boilerplate
This repo's name explain itself. Still have a lot of things todo

## Project Structure

```
.
├── api/                    # API specifications
│   └── openapi/           # OpenAPI/Swagger definitions
├── cmd/                   # Application entry points
│   ├── foo/              # Example service
│   └── svc-auth/         # Authentication service
│       ├── cmd/          # Service-specific commands
│       ├── config/       # Configuration management
│       ├── entity/       # Domain entities
│       ├── httpapi/      # HTTP API handlers
│       ├── migrations/   # Database migrations
│       ├── repository/   # Data access layer
│       └── service/      # Business logic layer
├── frontend/             # Frontend applications
│   └── auth/            # Authentication frontend (React/TypeScript)
├── pkg/                  # Shared Go packages
│   ├── conv/            # Conversion utilities
│   ├── core/            # Core abstractions (server, repository, service base)
│   ├── log/             # Logging utilities
│   ├── migrator/        # Database migration tools
│   └── tools/           # Build tools and dependencies
├── proto/               # Protocol Buffer definitions
│   ├── api/            # API service definitions
│   ├── common/         # Shared proto definitions
│   └── models/         # Data model definitions
├── rpc/                # Generated gRPC code (auto-generated from proto/)
├── scripts/            # Build and setup scripts
├── devenv/             # Development environment configuration
└── init/               # Initialization files
```

### Key Directories

- **`cmd/`**: Contains the main applications. Each subdirectory represents a separate microservice
- **`pkg/`**: Reusable Go packages shared across services
- **`proto/`**: Protocol Buffer definitions for gRPC services
- **`rpc/`**: Auto-generated Go code from protobuf definitions
- **`frontend/`**: Web frontend applications (React/TypeScript with Tailwind CSS)
- **`api/`**: API documentation and specifications

# Getting started
Before we start please make sure you have already installed these pieces of software
- Go >= 1.23.0
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
# create docker network
docker network create gpnetwork
# ramp up the environment
make docker.up
```

## Generate protobuf
You can define your protobuf inside `proto` directory and then run the following command to generate output
```bash
make generate-proto
```
All protobuf generated will be under `rpc` directory. To see how to implement API from generated code after run the command above, please refer to [Buf quick start](https://buf.build/docs/cli/quickstart/).

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

In Docker Desktop, go to Settings > Resources > WSL Integration, then disable and re-enable integration for your WSL2 distro.

Uncomment the `insecure: true` line in `devenv/traefik/traefik.yml` to enable insecure access to the Traefik dashboard.

