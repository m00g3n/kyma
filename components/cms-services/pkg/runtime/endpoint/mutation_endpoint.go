package endpoint

import (
	"context"
	"io"
	"net/http"

	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type mutationEndpoint struct {
	name    string
	mutator Mutator
}

// Mutator is the interface implemented by objects that can mutate objects
type Mutator interface {
	Mutate(ctx context.Context, reader io.Reader, metadata string) ([]byte, error)
}

var _ service.HTTPEndpoint = &mutationEndpoint{}

// NewMutation is the constructor that creates new Mutation Endpoint
func NewMutation(name string, mutator Mutator) service.HTTPEndpoint {
	return &mutationEndpoint{
		name:    name,
		mutator: mutator,
	}
}

// Name returns name of the endpoint
func (e *mutationEndpoint) Name() string {
	return e.name
}

// Handle process an HTTP request and calls mutator
func (e *mutationEndpoint) Handle(writer http.ResponseWriter, request *http.Request) {
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

	result, err := e.mutator.Mutate(request.Context(), content, metadata)
	if err != nil {
		log.Error(errors.Wrap(err, "while mutating request"))
		http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(result)
}
