FROM golang:1.25-alpine as build
WORKDIR /app
COPY . .

WORKDIR /app/service_A
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service_a ./cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/service_A .
ENTRYPOINT [ "./service_a" ]