package docstopic

import (
	"github.com/vrischmann/envconfig"
	"time"
)

type Config struct {
	DocsTopicRelistInterval time.Duration `envconfig:"default=5m"`
	BucketRegion            string        `envconfig:"optional"`
	WebHookCfgMapName       string        `envconfig:"default=webhook-config-map"`
	WebHookCfgMapNamespace  string        `envconfig:"default=kyma-system"`
}

func loadConfig(prefix string) (Config, error) {
	cfg := Config{}
	err := envconfig.InitWithPrefix(&cfg, prefix)
	return cfg, err
}
