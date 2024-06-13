FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o bin/reparttask ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/reparttask .

ENTRYPOINT ["./reparttask"]
