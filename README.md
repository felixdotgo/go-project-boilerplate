# Golang project boilerplate
This repo's name explain itself. Still have a lot of things todo

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
