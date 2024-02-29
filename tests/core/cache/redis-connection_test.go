package cache

import (
	"github.com/r4stl1n/micro-serv/pkg/core/cache"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

func TestRedisConnection(t *testing.T) {

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

	// Disconnect from redis
	disconnectError := redisClient.Disconnect()

	if disconnectError != nil {
		zap.L().Error("failed to disconnect from redis", zap.Error(disconnectError))
		t.Fatal()
	}
}
