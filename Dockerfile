# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o resolv ./cmd/resolv

# Final stage
FROM alpine:latest

# Install CA certificates to enable HTTPS requests to upstream DoH resolvers
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/resolv .

EXPOSE 53/udp

ENTRYPOINT ["./resolv", "-config", "config.toml"]
