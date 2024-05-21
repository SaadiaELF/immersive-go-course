# syntax=docker/dockerfile:1

# Builder stage
FROM golang:1.22-bullseye AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN mkdir -p out
RUN go build -o out/main main.go

# Final stage
FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/out/main /app/main

EXPOSE 80
CMD ["/app/main"]