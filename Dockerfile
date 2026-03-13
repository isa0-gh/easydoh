# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o easydoh ./cmd/easydoh

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/easydoh .

# Create config directory
RUN mkdir -p /etc/easydoh

EXPOSE 53/udp

ENTRYPOINT ["./easydoh"]
