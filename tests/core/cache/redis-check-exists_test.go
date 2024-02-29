package cache

import (
	"github.com/r4stl1n/micro-serv/pkg/core/cache"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestRedisCheckExists(t *testing.T) {

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

	// Insert an item into the cache item
	_, insertIntoCacheError := redisClient.AddCache("test", base.DummyMessage{Text: "herpderp"}, 5*time.Second)

	if insertIntoCacheError != nil {
		zap.L().Error("failed to insert cache item into redis", zap.Error(insertIntoCacheError))
		t.Fatal()
	}

	exists, cacheExistsError := redisClient.CacheExists("test")

	if cacheExistsError != nil {
		zap.L().Error("failed to check cache item in redis", zap.Error(cacheExistsError))
		t.Fatal()
	}

	if exists == false {
		zap.L().Error("failed to verify if item was in cache")
		t.Fatal()
	}

	// Delete the item from the cache
	deleteFromCacheError := redisClient.RemoveCache("test")

	if deleteFromCacheError != nil {
		zap.L().Error("failed to remove cache item from redis", zap.Error(deleteFromCacheError))
		t.Fatal()
	}

	// Disconnect from redis
	disconnectError := redisClient.Disconnect()

	if disconnectError != nil {
		zap.L().Error("failed to disconnect from redis", zap.Error(disconnectError))
		t.Fatal()
	}
}
