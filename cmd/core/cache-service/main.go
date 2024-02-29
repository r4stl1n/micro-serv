package main

import (
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"runtime"
)

// init create a new custom logger
func init() {
	new(util.Logger).Init()
}

func main() {

	// Create the cache service
	serviceMan, serviceManError := new(service.CacheService).Init("", 10)

	if serviceManError != nil {
		zap.L().Fatal("failed to create cache service", zap.Error(serviceManError))
	}

	// Start the cache service
	runError := serviceMan.Run()

	if runError != nil {
		zap.L().Fatal("failed to run cache service", zap.Error(runError))
	}

	runtime.Goexit()
}
