FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/ordersystem ./cmd/ordersystem

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/ordersystem /app/ordersystem
COPY sql/migrations /app/sql/migrations

EXPOSE 8000 8080 50051

CMD ["/app/ordersystem"]
