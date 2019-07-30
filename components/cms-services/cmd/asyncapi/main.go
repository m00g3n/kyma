package main

import (
	"context"
	"net/http"

	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/signal"
	"github.com/vrischmann/envconfig"

	"github.com/kyma-project/kyma/components/cms-services/pkg/endpoint/asyncapi"
	logpkg "github.com/kyma-project/kyma/components/cms-services/pkg/runtime/log"
	"github.com/kyma-project/kyma/components/cms-services/pkg/runtime/service"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type config struct {
	Verbose bool `envconfig:"default=true"`
	Service service.Config
}

func main() {
	cfg, err := loadConfig("APP")
	if err != nil {
		log.Fatal(errors.Wrap(err, "while loading configuration"))
	}

	logpkg.Setup(cfg.Verbose)

	stopCh := signal.SetupChannel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signal.CancelOnInterrupt(stopCh, ctx, cancel)

	srv := service.New(cfg.Service)

	log.Info("Registering endpoints")
	if err := asyncapi.AddToService(srv); err != nil {
		log.Fatal(errors.Wrap(err, "while registering endpoints"))
	}

	if err := srv.Start(ctx); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(errors.Wrap(err, "while starting service"))
		} else {
			log.Info("Service shutted down")
		}
	}
}

func loadConfig(prefix string) (config, error) {
	cfg := config{}
	err := envconfig.InitWithPrefix(&cfg, prefix)
	return cfg, err
}
