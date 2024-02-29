package main

import (
	"github.com/r4stl1n/micro-serv/internal/core/data-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"runtime"
)

// init create a new custom logger
func init() {
	new(util.Logger).Init()
}

func main() {

	// Create the data service
	serviceMan, serviceManError := new(service.DataService).Init("", "data")

	if serviceManError != nil {
		zap.L().Fatal("failed to create data service", zap.Error(serviceManError))
	}

	// Start the data service
	runError := serviceMan.Run()

	if runError != nil {
		zap.L().Fatal("failed to run data service", zap.Error(runError))
	}

	runtime.Goexit()
}
