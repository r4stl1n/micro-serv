package cache_service

import (
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/mediators"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCacheExists(t *testing.T) {

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

	// Add a item to the cache
	addCacheResponse, addCacheResponseError := cacheMediator.AddCache("test", 5*time.Second, base.DummyMessage{
		Text: "testData",
	})

	if addCacheResponseError != nil {
		zap.L().Error("failed to add cache", zap.Error(addCacheResponseError))
		t.Fatal()
	}

	// Attempt to retrieve a non-existent cache item
	cacheExistsResponse, cacheExistsResponseError := cacheMediator.CheckIfCacheExists("test")

	if cacheExistsResponseError != nil {
		zap.L().Error("failed to check cache exists", zap.Error(cacheExistsResponseError))
		t.Fatal()
	}

	if cacheExistsResponse.Exists == false {
		zap.L().Error("cache does not exists")
		t.Fatal()
	}

	// Remove an item from the cache
	_, removeCacheResponseError := cacheMediator.RemoveCache(addCacheResponse.Key)

	if removeCacheResponseError != nil {
		zap.L().Error("failed to remove cache", zap.Error(removeCacheResponseError))
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
