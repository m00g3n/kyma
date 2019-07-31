package endpoint

import (
	"context"
	"io"
	"net/http"

	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type validationEndpoint struct {
	name      string
	validator Validator
}

// Validator is the interface implemented by objects that can validate an request
type Validator interface {
	Validate(ctx context.Context, reader io.Reader, metadata string) error
}

var _ service.HTTPEndpoint = &validationEndpoint{}

// NewValidation is the constructor that creates new Validation Endpoint
func NewValidation(name string, validator Validator) service.HTTPEndpoint {
	return &validationEndpoint{
		name:      name,
		validator: validator,
	}
}

// Name returns name of the endpoint
func (e *validationEndpoint) Name() string {
	return e.name
}

// Handle process an HTTP request and calls validator
func (e *validationEndpoint) Handle(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := request.ParseMultipartForm(32 << 20); err != nil {
		log.Error(errors.Wrap(err, "while parsing multipart request"))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer request.MultipartForm.RemoveAll()

	content, _, err := request.FormFile("content")
	if err != nil {
		log.Error(errors.Wrap(err, "while accessing content"))
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer content.Close()

	metadata := request.FormValue("metadata")

	if err := e.validator.Validate(request.Context(), content, metadata); err != nil {
		log.Error(errors.Wrap(err, "while validating request"))
		http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
