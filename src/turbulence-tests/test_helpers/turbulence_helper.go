package test_helpers

import (
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/turbulence/client"
)

func TurbulenceClient() client.Turbulence {
	config := client.NewConfigFromEnv()
	clientLogger := logger.NewLogger(logger.LevelNone)
	return client.NewFactory(clientLogger).New(config)
}


