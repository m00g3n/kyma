package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"

	parser "github.com/asyncapi/parser/pkg"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/endpoint"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	"github.com/pkg/errors"
)

// AddValidation registers endpoint in service
func AddValidation(srv service.Service) error {
	validator := &validator{}

	srv.Register(endpoint.NewValidation("v1/validate", validator))
	return nil
}

var _ endpoint.Validator = &validator{}

type validator struct{}

func (v *validator) Validate(ctx context.Context, reader io.Reader, metadata string) error {
	document := v.streamToByte(reader)
	_, err := parser.Parse(document, false)
	if err != nil && len(err.ParsingErrors) > 0 {
		msg := err.ParsingErrors[0].String()
		for _, error := range err.ParsingErrors[1:] {
			msg = fmt.Sprintf("%s, %s", msg, error.String())
		}

		return errors.New(msg)
	}

	if err != nil {
		return err
	}

	return nil
}

func (v *validator) streamToByte(reader io.Reader) []byte {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)

	return buffer.Bytes()
}
