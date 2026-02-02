package main

import (
	"CloudRun/server/cmd/api"
	"CloudRun/server/cmd/internal/config"
	"CloudRun/server/cmd/internal/httpclient"
	"CloudRun/server/cmd/internal/service"
	"crypto/tls"
	"fmt"

	"net/http"
)

func main() {

	cfg := config.LoadConfig()
	httpClient := httpclient.New()
	service := service.NewWeatherService(cfg, httpClient)
	handler := api.NewWeatherHandler(service)

	// Desabilitar a verificação do certificado SSL
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", handler.GetWeatherByCityHandler)
	http.ListenAndServe(":8080", nil)

	fmt.Println("Server is running...")
}
