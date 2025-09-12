# syntax=docker/dockerfile:1.4

# Development stage with hot reload
FROM golang:1.24.2-alpine AS development
RUN apk add --no-cache git gcc musl-dev
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go workspace files first for better caching
COPY go.work go.work.sum ./

# Copy core module (required dependency)
COPY core/ ./core/

# Copy settings module
COPY settings/ ./settings/

WORKDIR /app/settings

# Download dependencies
RUN go mod download

# Expose settings service port
EXPOSE 8080

# Use air for hot reload in development
CMD ["air", "-c", ".air.toml", "-d"]

# Production stage with optimized binary
FROM golang:1.24.2-alpine AS builder
RUN apk update && apk add --no-cache git upx gcc musl-dev

WORKDIR /app

# Copy workspace configuration
COPY go.work go.work.sum ./

# Copy core module
COPY core/ ./core/

# Copy settings module
COPY settings/ ./settings/

WORKDIR /app/settings

# Cache dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

# Build optimized binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -v -o settings-service . && \
    upx --best --lzma settings-service

# Final production image
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Create non-root user
RUN addgroup -S settings && adduser -S settings -G settings

# Copy binary from builder
COPY --from=builder /app/settings/settings-service /app/settings-service
RUN chown settings:settings /app/settings-service

USER settings

EXPOSE 8080

CMD ["./settings-service"]