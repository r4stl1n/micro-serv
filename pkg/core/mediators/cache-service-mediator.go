package mediators

import (
	"fmt"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"time"
)

type CacheServiceMediator struct {
	natsClient *mq.Nats
	prefix     string
	version    consts.VersionConst
}

// Init creates a new cache mediator for interacting with the cache service
func (c *CacheServiceMediator) Init(nats *mq.Nats, version consts.VersionConst, prefix string) *CacheServiceMediator {
	*c = CacheServiceMediator{
		natsClient: nats,
		prefix:     prefix,
		version:    version,
	}

	return c
}

// AddCache adds data to a cache based on
func (c *CacheServiceMediator) AddCache(key string, expiration time.Duration, data any) (responses.AddCacheResponse, error) {

	var addCacheResponse responses.AddCacheResponse

	// Send the add cache request
	addCacheResponseError := c.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", c.version, c.prefix, consts.CoreCacheEndpointAddCache),
		requests.AddCacheRequest{
			Key:        key,
			Data:       data,
			Expiration: expiration,
		}, &addCacheResponse)

	return addCacheResponse, addCacheResponseError
}

// RemoveCache removes data from the cache
func (c *CacheServiceMediator) RemoveCache(key string) (responses.RemoveCacheResponse, error) {

	var removeCacheResponse responses.RemoveCacheResponse

	// Send a removal message
	removeCacheResponseError := c.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", c.version, c.prefix, consts.CoreCacheEndpointRemoveCache),
		requests.RemoveCacheRequest{
			Key: key,
		}, &removeCacheResponse)

	return removeCacheResponse, removeCacheResponseError
}

// CheckIfCacheExists checks if a key exists within the cache
func (c *CacheServiceMediator) CheckIfCacheExists(key string) (responses.CacheExistsResponse, error) {
	// Retrieve the cache
	var cacheExistsResponse responses.CacheExistsResponse

	cacheExistsResponseError := c.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", c.version, c.prefix, consts.CoreCacheEndpointCacheExists),
		requests.RetrieveCacheRequest{
			Key: key,
		}, &cacheExistsResponse)

	return cacheExistsResponse, cacheExistsResponseError
}

// RetrieveCache attempts to retrieve and decode data from the cache
func (c *CacheServiceMediator) RetrieveCache(key string, output any) error {
	// Retrieve the cache
	var retrieveCacheResponse responses.RetrieveCacheResponse

	retrieveCacheResponseError := c.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", c.version, c.prefix, consts.CoreCacheEndpointRetrieveCache),
		requests.RetrieveCacheRequest{
			Key: key,
		}, &retrieveCacheResponse)

	if retrieveCacheResponseError != nil {
		return retrieveCacheResponseError
	}

	// Attempt to decode the message
	return retrieveCacheResponse.DecodeData(output)

}

// RequestKeyLock attempts to request a key lock on an id
func (c *CacheServiceMediator) RequestKeyLock(key string, expiration time.Duration) (responses.RequestKeyLockResponse, error) {

	var requestKeyLockResponse responses.RequestKeyLockResponse

	// Send the add cache request
	requestKeyLockResponseError := c.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", c.version, c.prefix, consts.CoreCacheEndpointRequestKeyLock),
		requests.RequestKeyLockRequest{
			Key:        key,
			Expiration: expiration,
		}, &requestKeyLockResponse)

	return requestKeyLockResponse, requestKeyLockResponseError
}

// RemoveKeyLock removes an existing key lock
func (c *CacheServiceMediator) RemoveKeyLock(key string) (responses.RemoveKeyLockResponse, error) {

	var removeKeyLockResponse responses.RemoveKeyLockResponse

	// Send the add cache request
	removeKeyLockResponseError := c.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", c.version, c.prefix, consts.CoreCacheEndpointRemoveKeyLock),
		requests.RemoveKeyLockRequest{
			Key: key,
		}, &removeKeyLockResponse)

	return removeKeyLockResponse, removeKeyLockResponseError

}
