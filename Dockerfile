# 1. Build stage: Compile the Go program
ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION}-alpine AS builder

# Install required tools (e.g., certificates, Git)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy module files first (use caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code and compile
COPY . .
RUN go build -o oauth2-server ./cmd/server

# 2. Runtime stage: Minimal container
FROM alpine:latest

# Install CA certificates if TLS is used (e.g., for JWKS over HTTPS)
RUN apk add --no-cache ca-certificates

# Copy binary - no keys
COPY --from=builder /app/oauth2-server /oauth2-server

# Runtime expects /app/keys – will be mounted
VOLUME ["/app/keys"]

EXPOSE 8080

# Define entrypoint
ENTRYPOINT ["/oauth2-server"]