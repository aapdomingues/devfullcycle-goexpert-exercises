FROM golang:1.25 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun server/cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
ENTRYPOINT [ "./cloudrun" ]