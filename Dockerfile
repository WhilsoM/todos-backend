FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-api ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/todo-api ./todo-api
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./todo-api"]
