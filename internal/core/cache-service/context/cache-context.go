package context

import (
	"github.com/r4stl1n/micro-serv/pkg/core/cache"
)

type CacheContext struct {
	Cache *cache.Redis
}

// Init create the scaffold context
func (c *CacheContext) Init(cache *cache.Redis) *CacheContext {

	*c = CacheContext{
		Cache: cache,
	}

	return c
}
