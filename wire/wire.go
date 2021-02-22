//+build wireinject

package main

import (
	"learnGo/wire/config"
	"learnGo/wire/repo"
	"learnGo/wire/service"

	"github.com/google/wire"
)

func initService() (service.Service, error) {
	wire.Build(
		config.NewConfig,
		repo.NewRepo,
		service.NewService,
	)
	return nil, nil
}
