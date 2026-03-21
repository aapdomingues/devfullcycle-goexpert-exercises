FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/ordersystem ./cmd/ordersystem

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/ordersystem /app/ordersystem
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY sql/migrations /app/sql/migrations

EXPOSE 8000 8080 50051
