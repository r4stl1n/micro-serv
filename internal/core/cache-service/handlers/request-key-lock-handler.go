package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/cache-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func RequestKeyLockHandler(cacheContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	cContext := cacheContext.(*context.CacheContext)

	// Extract keyLockRequest message
	var keyLockRequest requests.RequestKeyLockRequest

	decodeError := sContext.DecodeRequest(microReq, &keyLockRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Request a key lock
	acquired, requestKeyLockError := cContext.Cache.RequestKeyLock(keyLockRequest.Key, keyLockRequest.Expiration)

	if requestKeyLockError != nil {
		sContext.PublishError(microReq, requestKeyLockError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.RequestKeyLockResponse{Key: keyLockRequest.Key, Acquired: acquired})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
