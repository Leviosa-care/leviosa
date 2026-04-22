# syntax=docker/dockerfile:1.4

FROM golang:1.24.2-alpine AS builder
RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o leviosa ./cmd/app/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

COPY --from=builder /app/leviosa /app/leviosa

# Change ownership of /app to appuser
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 3500 3501

CMD ["./leviosa"]
