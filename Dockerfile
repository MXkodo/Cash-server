FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y wget && \
    apt-get clean

RUN wget https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz && \
    tar xzf migrate.linux-amd64.tar.gz && \
    mv migrate /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate && \
    rm migrate.linux-amd64.tar.gz

WORKDIR /app

COPY --from=builder /app/main .

COPY internal/repo/postgresql/migrations /app/migrations

ENV PATH="/usr/local/bin:${PATH}"

CMD ["./main"]
