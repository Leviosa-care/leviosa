FROM golang:1.24.1-alpine AS builder
RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s" -v -o leviosa ./cmd/leviosa/*.go

FROM alpine:latest

RUN apk add --no-cache sqlite

WORKDIR /app

COPY --from=builder /app/leviosa .
COPY --from=builder /app/migrations/production ./migrations/production/
COPY --from=builder /app/*.env ./

RUN adduser -D runner

RUN mkdir -p /app/data && chown -R runner:runner /app

USER runner

EXPOSE 3500 3501
CMD ["./leviosa"]
