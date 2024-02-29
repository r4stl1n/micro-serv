package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func RetrieveCacheHandler(cacheContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	cContext := cacheContext.(*context.CacheContext)

	// Extract retrieveCacheRequest message
	var retrieveCacheRequest requests.RetrieveCacheRequest

	decodeError := sContext.DecodeRequest(microReq, &retrieveCacheRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Retrieve the data from the cache
	// We use an interface here since we support all data
	// That can be serialized to json string
	var data any
	retrieveError := cContext.Cache.RetrieveCache(retrieveCacheRequest.Key, &data)

	if retrieveError != nil {
		sContext.PublishError(microReq, retrieveError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.RetrieveCacheResponse{Key: retrieveCacheRequest.Key, Data: data})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
