package cache_service

import (
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/mediators"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestRequestKeyLockFail(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Create the service
	serviceMan, serviceManError := new(service.CacheService).Init(consts.TestPrefixConst.ToString(), 0)

	if serviceManError != nil {
		zap.L().Fatal("failed to create cache service", zap.Error(serviceManError))
	}

	// Start the service
	runError := serviceMan.Run()

	if runError != nil {
		zap.L().Fatal("failed to run data service", zap.Error(runError))
	}

	// Create a new nats client from the environment variables
	natsClient, natsClientError := new(mq.Nats).Init()

	if natsClientError != nil {
		t.Skip("nats not configured for test")
	}

	// Connect to the nats server
	connectError := natsClient.Connect()

	if connectError != nil {
		zap.L().Error("failed to connect to nats server")
		t.Fatal()
	}

	// Create a new cache mediator
	cacheMediator := new(mediators.CacheServiceMediator).Init(natsClient, consts.VersionV1Const, consts.TestPrefixConst.ToString())

	// Send the add cache request
	requestKeyLockResponse, requestKeyLockResponseError := cacheMediator.RequestKeyLock("testL", 5*time.Second)

	if requestKeyLockResponseError != nil {
		zap.L().Error("failed to retrieve key lock response from cache", zap.Error(requestKeyLockResponseError))
		t.Fatal()
	}

	if requestKeyLockResponse.Acquired == false {
		zap.L().Error("failed to acquire lock")
		t.Fatal()
	}

	// Send the add cache request
	requestKeyLockResponse2, requestKeyLockResponseError2 := cacheMediator.RequestKeyLock("testL", 5*time.Second)

	if requestKeyLockResponseError2 != nil {
		zap.L().Error("failed to retrieve key lock response from cache", zap.Error(requestKeyLockResponseError2))
		t.Fatal()
	}

	if requestKeyLockResponse2.Acquired == true {
		zap.L().Error("acquired lock we should not have been able to")
		t.Fatal()
	}

	// Send the add cache request
	_, removeKeyLockResponseError := cacheMediator.RemoveKeyLock(requestKeyLockResponse.Key)

	if removeKeyLockResponseError != nil {
		zap.L().Error("failed to retrieve remove key lock response from cache", zap.Error(removeKeyLockResponseError))
		t.Fatal()
	}
	// Shut it all down
	serviceStopError := serviceMan.Stop()

	if serviceStopError != nil {
		zap.L().Error("failed to stop the service", zap.Error(serviceStopError))
		t.Fatal()
	}

	natsClient.Disconnect()

}
