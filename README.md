# Golang project boilerplate
This repo's name explain itself. Still have a lot of things todo

# Getting started
Before we start please make sure you're already installed these pieces of software
- Go >= 1.21.0
- protoc >= 25.1
- buf >= 1.47.2
- Docker
- Cmake/Make

## Create a new service
To create a new service, just simply make a copy of `cmd/foo` directory and rename it to whatever you want

## Generate protobuf
You can define your protobuf inside `proto` directory and then run the following command to generate output
```bash
make generate-proto
```
All protobuf generated will be under `rpc` directory. To see how to implement API from generated code after run the command above, please refer to [Buf quick start](https://buf.build/docs/cli/quickstart/).
