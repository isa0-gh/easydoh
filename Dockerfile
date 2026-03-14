# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o resolv ./cmd/resolv


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/resolv /usr/local/bin/resolv

EXPOSE 53/udp

ENTRYPOINT ["resolv", "-config", "/app/config.toml"]
