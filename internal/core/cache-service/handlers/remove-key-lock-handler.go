package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func RemoveKeyLockHandler(cacheContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	cContext := cacheContext.(*context.CacheContext)

	// Extract keyLockRequest message
	var keyLockRequest requests.RemoveKeyLockRequest

	decodeError := sContext.DecodeRequest(microReq, &keyLockRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Request a key lock
	removeKeyLockError := cContext.Cache.RemoveKeyLock(keyLockRequest.Key)

	if removeKeyLockError != nil {
		sContext.PublishError(microReq, removeKeyLockError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.RemoveKeyLockResponse{Key: keyLockRequest.Key})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
