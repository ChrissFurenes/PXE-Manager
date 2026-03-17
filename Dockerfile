# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /pxe-manager ./cmd/pxe-server

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libc6 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /pxe-manager /usr/local/bin/pxe-manager
COPY config.yaml /app/config.yaml
COPY pxe /app/pxe
COPY web /app/web

EXPOSE 69/udp
EXPOSE 8080/tcp

CMD ["/usr/local/bin/pxe-manager"]