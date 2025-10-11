FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY ./pkg/ ./pkg
COPY ./rpc/ ./rpc
COPY cmd/foo ./cmd/foo

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/foo/main.go


FROM alpine:3.21 AS run

WORKDIR /app
COPY --from=build /server ./server
EXPOSE 8000

CMD ["/app/server"]
