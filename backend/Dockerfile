FROM golang:1.21 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 go build -o main ./cmd/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]