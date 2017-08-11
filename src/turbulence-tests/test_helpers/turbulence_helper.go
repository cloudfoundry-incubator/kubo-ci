package test_helpers

import (
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/turbulence/client"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/tasks"
)

func TurbulenceClient() client.Turbulence {
	config := client.NewConfigFromEnv()
	clientLogger := logger.NewLogger(logger.LevelNone)
	return client.NewFactory(clientLogger).New(config)
}


func KillVms(vmType, limitString string) client.Incident {
	turbulenceClient := TurbulenceClient()
	limit, _ := selector.NewLimitFromString(limitString)
	req := incident.Request{
		Tasks: tasks.OptionsSlice{tasks.KillOptions{}},
		Selector: selector.Request{
			Deployment: &selector.NameRequest{Name: "ci-service"},
			Group:      &selector.NameRequest{Name: vmType},
			ID:         &selector.IDRequest{Limit: limit},
		},
	}
	inc := turbulenceClient.CreateIncident(req)
	inc.Wait()
	return inc
}
