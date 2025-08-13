# Этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка под Linux amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./core/cmd/main.go

# Этап выполнения
FROM alpine:3.22
RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /app/main /usr/local/bin/main

EXPOSE 8082
ENTRYPOINT ["/usr/local/bin/main"]
