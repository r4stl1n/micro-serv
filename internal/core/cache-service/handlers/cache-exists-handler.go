package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func CacheExistsHandler(cacheContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	cContext := cacheContext.(*context.CacheContext)

	// Extract cacheExistsRequest message
	var cacheExistsRequest requests.CacheExistsRequest

	decodeError := sContext.DecodeRequest(microReq, &cacheExistsRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Attempt to check if the cache exists
	exists, existsError := cContext.Cache.CacheExists(cacheExistsRequest.Key)

	if existsError != nil {
		sContext.PublishError(microReq, existsError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.CacheExistsResponse{Key: cacheExistsRequest.Key, Exists: exists})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
