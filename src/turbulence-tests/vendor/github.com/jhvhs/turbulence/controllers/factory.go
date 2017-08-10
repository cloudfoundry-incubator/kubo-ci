package controllers

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/jhvhs/turbulence/incident"
	"github.com/jhvhs/turbulence/scheduledinc"
	"github.com/jhvhs/turbulence/tasks"
)

type FactoryRepos interface {
	IncidentsRepo() incident.Repo
	ScheduledIncidentsRepo() scheduledinc.Repo
	TasksRepo() tasks.Repo
}

type Factory struct {
	HomeController               HomeController
	IncidentsController          IncidentsController
	ScheduledIncidentsController ScheduledIncidentsController
	TasksController              TasksController
}

func NewFactory(r FactoryRepos, logger boshlog.Logger) (Factory, error) {
	isRepo := r.IncidentsRepo()
	sisRepo := r.ScheduledIncidentsRepo()
	arRepo := r.TasksRepo()

	factory := Factory{
		HomeController:               NewHomeController(isRepo, sisRepo, logger),
		IncidentsController:          NewIncidentsController(isRepo, logger),
		ScheduledIncidentsController: NewScheduledIncidentsController(sisRepo, logger),
		TasksController:              NewTasksController(arRepo, logger),
	}

	return factory, nil
}
