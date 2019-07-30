package v1

import (
	v2 "github.com/asyncapi/converter-go/pkg/converter/v2"

	"github.com/asyncapi/converter-go/pkg/decode"
	"github.com/asyncapi/converter-go/pkg/encode"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/endpoint"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"

	"bufio"
	"bytes"
	"context"
	"io"
)

type Convert func(reader io.Reader, writer io.Writer) error

var _ endpoint.Mutator = Convert(nil)

func (c Convert) Mutate(ctx context.Context, contentType string, reader io.Reader, metadata string) ([]byte, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err := c(reader, writer); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func AddConversion(srv service.Service) error {
	converter, err := v2.New(decode.FromJSONWithYamlFallback, encode.ToJSON)
	if err != nil {
		return nil
	}
	convert := Convert(converter.Convert)
	srv.Register(endpoint.NewMutation("v1/convert", convert))
	return nil
}
