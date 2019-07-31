package fake

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	log "github.com/sirupsen/logrus"
)

// Service is a fake implementation of Asset Store service
type Service struct {
	endpoints []service.HTTPEndpoint
	mux       *http.ServeMux
}

var _ service.Service = &Service{}

// NewService is a constructor that creates new fake service
func NewService() *Service {
	return &Service{}
}

// RequestBodyFromFile builds multipart request from file
func RequestBodyFromFile(filePath, parameters string) (io.Reader, string, error) {
	buffer := &bytes.Buffer{}
	formWriter := multipart.NewWriter(buffer)
	defer formWriter.Close()

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, "", errors.Wrapf(err, "while opening file %s", filePath)
		}
		defer file.Close()

		contentWriter, err := formWriter.CreateFormFile("content", filepath.Base(file.Name()))
		if err != nil {
			return nil, "", errors.Wrapf(err, "while creating content field for file %s", filePath)
		}

		_, err = io.Copy(contentWriter, file)
		if err != nil {
			return nil, "", errors.Wrapf(err, "while copying file %s to content field", filePath)
		}
	}

	if parameters != "" {
		err := formWriter.WriteField("parameters", parameters)
		if err != nil {
			return nil, "", errors.Wrapf(err, "while creating parameters field for parameters %s", parameters)
		}
	}
	return buffer, formWriter.FormDataContentType(), nil
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (s *Service) ServeHTTP(method, endpoint, contentType string, body io.Reader) *http.Response {
	recorder := httptest.NewRecorder()
	if s.mux == nil {
		http.Error(recorder, "Server is not initialized", http.StatusInternalServerError)
	}

	request := httptest.NewRequest(method, endpoint, body)
	request.Header.Add("Content-Type", contentType)

	s.mux.ServeHTTP(recorder, request)
	return recorder.Result()
}

// Start configure routes in fake service
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

// Register adds endpoints to service
func (s *Service) Register(endpoint service.HTTPEndpoint) {
	s.endpoints = append(s.endpoints, endpoint)
}
