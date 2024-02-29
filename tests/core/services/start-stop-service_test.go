package services

import (
	"github.com/r4stl1n/micro-serv/internal/core/dummy-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

func TestStartStopService(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Being lazy and using this to check if nats is running
	_, natsClientError := new(mq.Nats).Init()

	if natsClientError != nil {
		t.Skip("nats not configured for test")
	}

	// Create the dummy service
	dummyService, dummyServiceError := new(service.DummyService).Init("dummy")

	if dummyServiceError != nil {
		zap.L().Error("failed to create dummy service", zap.Error(dummyServiceError))
		t.Fatal()
	}

	// Run the dummy service as a separate go function
	runError := dummyService.Run()

	if runError != nil {
		zap.L().Error("failed to run dummy service", zap.Error(runError))
	}

	// Stop the service
	stopError := dummyService.Stop()

	if stopError != nil {
		zap.L().Error("failed to stop dummy service", zap.Error(stopError))
		t.Fatal()
	}

	if dummyService.IsRunning() == true {
		zap.L().Fatal("service did not stop")
	}
}
