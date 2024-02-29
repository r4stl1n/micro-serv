package cache

import (
	"github.com/r4stl1n/micro-serv/pkg/core/cache"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

func TestRedisRetrieveFail(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Being lazy and using this to check if nats is running
	redisClient, redisClientError := new(cache.Redis).Init(0)

	if redisClientError != nil {
		t.Skip("redis not configured for test")
	}

	// Connect to redis
	connectError := redisClient.Connect()

	if connectError != nil {
		zap.L().Error("failed to connect to redis", zap.Error(connectError))
		t.Fatal()
	}

	// Retrieve data from cache
	var result base.DummyMessage
	retrieveError := redisClient.RetrieveCache("herpderp123", &result)

	if retrieveError == nil {
		zap.L().Error("retrieve an item and we should not have")
		t.Fatal()
	}

}
