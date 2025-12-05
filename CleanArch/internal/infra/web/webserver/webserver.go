package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router         chi.Router
	HandlersConfig []HandlerConfig
	WebServerPort  string
}

type HandlerConfig struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func NewHandlerConfig(method, path string, handler http.HandlerFunc) HandlerConfig {
	return HandlerConfig{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:         chi.NewRouter(),
		HandlersConfig: []HandlerConfig{},
		WebServerPort:  serverPort,
	}
}

func (s *WebServer) AddHandler(handlerConfig HandlerConfig) {
	s.HandlersConfig = append(s.HandlersConfig, handlerConfig)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, handlerConfig := range s.HandlersConfig {
		s.Router.MethodFunc(handlerConfig.Method, handlerConfig.Path, handlerConfig.Handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
