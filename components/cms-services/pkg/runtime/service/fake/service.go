package fake

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	endpoints []service.HttpEndpoint
	mux       *http.ServeMux
	recorder  *httptest.ResponseRecorder
}

var _ service.Service = &Service{}

func NewService() *Service {
	return &Service{
		recorder: httptest.NewRecorder(),
	}
}

func (s *Service) ServeHTTP(method, endpoint, contentType string, body io.Reader) *http.Response {
	if s.mux == nil {
		http.Error(s.recorder, "Server is not initialized", http.StatusInternalServerError)
	}

	request := httptest.NewRequest(method, endpoint, body)
	request.Header.Add("Contnet-Type", contentType)

	s.mux.ServeHTTP(s.recorder, request)
	return s.recorder.Result()
}

func (s *Service) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	for _, endpoint := range s.endpoints {
		log.Infof("Registering %s endpoint", endpoint.Name())
		path := fmt.Sprintf("/%s", endpoint.Name())
		mux.HandleFunc(path, endpoint.Handle)
	}

	s.mux = mux
	return nil
}

func (s *Service) Register(endpoint service.HttpEndpoint) {
	s.endpoints = append(s.endpoints, endpoint)
}
