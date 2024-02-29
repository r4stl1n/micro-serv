package service

import (
	"fmt"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/handlers"
	"github.com/r4stl1n/micro-serv/pkg/core/cache"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
)

type CacheService struct {
	scaffold       *scaff.ServiceScaffold
	cacheCtx       *context.CacheContext
	endpointPrefix string
}

// Init create a new DataService
func (c *CacheService) Init(endpointPrefix string, cacheDb int) (*CacheService, error) {

	// Create a new scaffold instance
	scaffold, scaffoldingError := new(scaff.ServiceScaffold).Init("cache-service", "1.0.0", "Cache Service")

	if scaffoldingError != nil {
		return nil, scaffoldingError
	}

	// Create cache connection
	redisCache, redisCacheError := new(cache.Redis).Init(cacheDb)

	if redisCacheError != nil {
		return nil, redisCacheError
	}

	// Connect to cache
	cacheConnectError := redisCache.Connect()
	if cacheConnectError != nil {
		return nil, cacheConnectError
	}

	*c = CacheService{
		cacheCtx:       new(context.CacheContext).Init(redisCache),
		scaffold:       scaffold,
		endpointPrefix: endpointPrefix,
	}
	return c, nil
}

// StartUp create the handlers and perform any additional first time needs
func (c *CacheService) addHandlers() error {

	addCacheHandlerError := c.scaffold.AddHandler(c.cacheCtx,
		fmt.Sprintf("%s%s", c.endpointPrefix, consts.CoreCacheEndpointAddCache), handlers.AddCacheHandler)

	if addCacheHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			c.endpointPrefix, consts.CoreCacheEndpointAddCache, addCacheHandlerError)
	}

	removeCacheHandlerError := c.scaffold.AddHandler(c.cacheCtx,
		fmt.Sprintf("%s%s", c.endpointPrefix, consts.CoreCacheEndpointRemoveCache), handlers.RemoveCacheHandler)

	if removeCacheHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			c.endpointPrefix, consts.CoreCacheEndpointRemoveCache, removeCacheHandlerError)
	}

	retrieveCacheHandlerError := c.scaffold.AddHandler(c.cacheCtx,
		fmt.Sprintf("%s%s", c.endpointPrefix, consts.CoreCacheEndpointRetrieveCache), handlers.RetrieveCacheHandler)

	if retrieveCacheHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			c.endpointPrefix, consts.CoreCacheEndpointRetrieveCache, retrieveCacheHandlerError)
	}

	cacheExistsHandlerError := c.scaffold.AddHandler(c.cacheCtx,
		fmt.Sprintf("%s%s", c.endpointPrefix, consts.CoreCacheEndpointCacheExists), handlers.CacheExistsHandler)

	if cacheExistsHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			c.endpointPrefix, consts.CoreCacheEndpointCacheExists, cacheExistsHandlerError)
	}

	requestKeyLockHandler := c.scaffold.AddHandler(c.cacheCtx,
		fmt.Sprintf("%s%s", c.endpointPrefix, consts.CoreCacheEndpointRequestKeyLock), handlers.RequestKeyLockHandler)

	if requestKeyLockHandler != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			c.endpointPrefix, consts.CoreCacheEndpointRequestKeyLock, requestKeyLockHandler)
	}

	removeKeyLockHandler := c.scaffold.AddHandler(c.cacheCtx,
		fmt.Sprintf("%s%s", c.endpointPrefix, consts.CoreCacheEndpointRemoveKeyLock), handlers.RemoveKeyLockHandler)

	if removeKeyLockHandler != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			c.endpointPrefix, consts.CoreCacheEndpointRemoveKeyLock, removeKeyLockHandler)
	}

	return nil
}

// Run the service
func (c *CacheService) Run() error {

	startError := c.scaffold.Start()

	if startError != nil {
		return fmt.Errorf("failed to start the service: %s", startError)
	}

	return c.addHandlers()
}

// Stop the service
func (c *CacheService) Stop() error {
	stopError := c.scaffold.Stop()

	if stopError != nil {
		return fmt.Errorf("failed to stop the service: %s", stopError)
	}

	return nil
}

// IsRunning returns if the service is running
func (c *CacheService) IsRunning() bool {
	return c.scaffold.IsRunning()
}
