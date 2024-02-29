package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func AddCacheHandler(cacheContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	cContext := cacheContext.(*context.CacheContext)

	// Extract addCacheRequest message
	var addCacheRequest requests.AddCacheRequest

	decodeError := sContext.DecodeRequest(microReq, &addCacheRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Add data to cache
	addId, addError := cContext.Cache.AddCache(addCacheRequest.Key, addCacheRequest.Data, addCacheRequest.Expiration)

	if addError != nil {
		sContext.PublishError(microReq, addError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.AddCacheResponse{Key: addId})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
