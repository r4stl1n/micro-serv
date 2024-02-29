package main

import (
	"github.com/r4stl1n/micro-serv/internal/core/dummy-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"runtime"
)

// init create a new custom logger
func init() {
	new(util.Logger).Init()
}

func main() {

	// Create the dummy service
	serviceMan, serviceManError := new(service.DummyService).Init("")

	if serviceManError != nil {
		zap.L().Fatal("failed to create dummy service", zap.Error(serviceManError))
	}

	// Start the dummy service
	runError := serviceMan.Run()

	if runError != nil {
		zap.L().Fatal("failed to run dummy service", zap.Error(runError))
	}

	runtime.Goexit()
}
