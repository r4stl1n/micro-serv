package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func RemoveCacheHandler(cacheContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	cContext := cacheContext.(*context.CacheContext)

	// Extract removeCacheRequest message
	var removeCacheRequest requests.RemoveCacheRequest

	decodeError := sContext.DecodeRequest(microReq, &removeCacheRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Remove data from cache
	deleteError := cContext.Cache.RemoveCache(removeCacheRequest.Key)

	if deleteError != nil {
		sContext.PublishError(microReq, deleteError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.RemoveCacheResponse{Key: removeCacheRequest.Key})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
