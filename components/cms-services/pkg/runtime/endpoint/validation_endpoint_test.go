package endpoint_test

import (
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/endpoint"
	"github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidationEndpoint_Handle(t *testing.T) {
	for testName, testCase := range map[string]struct {
		targetMethod   string
		expectedStatus int
		body           io.Reader
		contentType    string
		validator      endpoint.Validator
	}{
		//"OK": {
		//	expectedStatus: http.StatusNotFound,
		//	targetEndpoint: "/test",
		//	targetMethod:   http.MethodPost,
		//},
		"invalid method": {
			expectedStatus: http.StatusMethodNotAllowed,
			targetMethod:   http.MethodGet,
		},
		"invalid request": {
			expectedStatus: http.StatusBadRequest,
			targetMethod:   http.MethodPost,
			body:           strings.NewReader("test"),
		},
		"missing file": {
			expectedStatus: http.StatusBadRequest,
			targetMethod:   http.MethodPost,
			body:           strings.NewReader("------SPLIT\r\nContent-Disposition: form-data; name=\"content\"; filename=\"test.ok.yaml\"\r\nContent-Type: text/yaml\r\n\r\n\r\n------SPLIT\r\nContent-Disposition: form-data; name=\"metadata\"\r\n\r\n\r\n------SPLIT--"),
			contentType:    "multipart/form-data; boundary=----SPLIT",
		},
		"validation failed": {
			expectedStatus: http.StatusUnprocessableEntity,
			targetMethod:   http.MethodPost,
		},
	} {
		t.Run(testName, func(t *testing.T) {
			// given
			g := gomega.NewWithT(t)
			edp := endpoint.NewValidation("test", nil)

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(edp.Handle)
			request := httptest.NewRequest(testCase.targetMethod, "/test", testCase.body)
			request.Header.Add("content-type", testCase.contentType)

			// when
			handler.ServeHTTP(recorder, request)

			// then
			g.Expect(recorder.Result().StatusCode).To(gomega.Equal(testCase.expectedStatus))
		})
	}
}
