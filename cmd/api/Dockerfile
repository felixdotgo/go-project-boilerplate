FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./pkg/ ./pkg
COPY ./rpc/ ./rpc
COPY cmd/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api/main.go


FROM alpine:3.21 AS run

WORKDIR /app
COPY --from=build /server ./server
EXPOSE 8000

CMD ["/app/server"]
