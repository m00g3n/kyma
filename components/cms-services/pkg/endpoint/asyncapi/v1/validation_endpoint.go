package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/asyncapi/parser/pkg"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/endpoint"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	"github.com/pkg/errors"
)

func AddValidation(srv service.Service) error {
	validator := &Validator{}

	srv.Register(endpoint.NewValidation("v1/validate", validator))
	return nil
}

var _ endpoint.Validator = &Validator{}

type Validator struct{}

func (v *Validator) Validate(ctx context.Context, contentType string, reader io.Reader, metadata string) error {
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

func (v *Validator) streamToByte(reader io.Reader) []byte {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)

	return buffer.Bytes()
}
