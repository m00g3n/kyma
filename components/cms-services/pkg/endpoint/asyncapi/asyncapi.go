package asyncapi

import (
	v1 "github.com/kyma-project/kyma/components/cms-services/pkg/endpoint/asyncapi/v1"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
)

// AddToServiceFuncs contains collection of functions that registers endpoints
var AddToServiceFuncs []func(service.Service) error

func init() {
	AddToServiceFuncs = append(AddToServiceFuncs, v1.AddValidation, v1.AddConversion)
}

// AddToService registers endpoints in service
func AddToService(s service.Service) error {
	for _, f := range AddToServiceFuncs {
		if err := f(s); err != nil {
			return err
		}
	}
	return nil
}
