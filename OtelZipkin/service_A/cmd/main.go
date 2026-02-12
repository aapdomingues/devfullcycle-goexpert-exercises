package main

import (
	"OtelZipkin/platform/httpclient"
	"OtelZipkin/platform/otel"
	"OtelZipkin/service_A/cmd/api"
	"OtelZipkin/service_A/cmd/config"
	"OtelZipkin/service_A/cmd/service"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	cfg := config.LoadConfig()
	httpClient := httpclient.New()
	service := service.NewWeatherService(cfg, httpClient)

	handler := api.NewWeatherHandler(service)

	// Desabilitar a verificação do certificado SSL
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.Handle(
		"/",
		otelhttp.NewHandler(
			http.HandlerFunc(handler.GetWeatherByCityHandler),
			"POST /weather",
		),
	)

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("Server is running...")
	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
