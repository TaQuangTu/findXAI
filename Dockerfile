# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install protoc and dependencies
RUN apk add --no-cache \
    protoc \
    protobuf-dev \
    make \
    git

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

# Runtime stage
FROM alpine:3.21.3

WORKDIR /app

# Copy the binary and protobuf definitions
COPY --from=builder /server /app/server
COPY --from=builder /app/pkg/protogen /app/pkg/protogen

# Install grpc_health_probe
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Install CA certificates
RUN apk add --no-cache ca-certificates

EXPOSE 50051

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["grpc_health_probe", "-addr=:50051"]

ENTRYPOINT ["/app/server"]