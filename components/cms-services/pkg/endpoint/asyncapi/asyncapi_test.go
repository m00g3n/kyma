package asyncapi_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/kyma-project/kyma/components/cms-services/pkg/endpoint/asyncapi"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service/fake"
	"github.com/onsi/gomega"
)

func TestV1Validation(t *testing.T) {
	g := gomega.NewWithT(t)

	srv, err := initService()
	g.Expect(err).ToNot(gomega.HaveOccurred())

	response := serveValidate(srv, "")
	g.Expect(response.StatusCode).To(gomega.Equal(http.StatusOK))
}

func serveValidate(srv *fake.Service, filePath string) *http.Response {
	response := srv.ServeHTTP(http.MethodPost, "/v1/validate", "", nil)

	return response
}

func buildQuery(filePath string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	w.CreateFormFile("content", filePath)
}

func initService() (*fake.Service, error) {
	srv := fake.NewService()
	asyncapi.AddToService(srv)

	srv.Start(context.TODO())

	return srv, nil
}
