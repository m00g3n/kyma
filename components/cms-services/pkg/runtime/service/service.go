package service

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Host string `envconfig:"default=127.0.0.1"`
	Port int    `envconfig:"default=3000"`
}

type Service interface {
	Register(endpoint HttpEndpoint)
	Start(ctx context.Context) error
}

type HttpEndpoint interface {
	Name() string
	Handle(writer http.ResponseWriter, request *http.Request)
}

type service struct {
	endpoints []HttpEndpoint
	host      string
	port      int
}

var _ Service = &service{}

func New(config Config) *service {
	return &service{
		host: config.Host,
		port: config.Port,
	}
}

func (s *service) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	for _, endpoint := range s.endpoints {
		log.Infof("Registering %s endpoint", endpoint.Name())
		path := fmt.Sprintf("/%s", endpoint.Name())
		mux.HandleFunc(path, endpoint.Handle)
	}

	host := fmt.Sprintf("%s:%d", s.host, s.port)

	srv := &http.Server{Addr: host, Handler: mux}
	log.Infof("Service listen at %s", host)

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Errorf("HTTP server Shutdown: %v", err)
		}
	}()

	return srv.ListenAndServe()
}

func (s *service) Register(endpoint HttpEndpoint) {
	s.endpoints = append(s.endpoints, endpoint)
}
