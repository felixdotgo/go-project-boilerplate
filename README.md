# Golang project boilerplate
This repo's name explain itself. Still have a lot of things todo

# Get started
Before we start please make sure you're already installed these pieces of software
- Go >= 1.21.0
- protoc >= 25.1
- buf >= 1.47.2
- Docker
- Cmake/Make

## Migration commands
Create migration in `migrations/sql` directory
```
go run cmd/migrator/migrator.go create -n my_file_name
```
Execute all migrations
```
go run cmd/migrator/migrator.go up
```

## Generate protobuf
You can define your protobuf inside `proto` directory and then run the following command to generate output
```bash
buf generate --include-imports
```
All protobuf generated will be under `rpc` directory.
