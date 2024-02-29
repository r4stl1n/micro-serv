package scaff

import (
	"encoding/json"
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"go.uber.org/zap"
)

type ScaffoldContext struct {
	Nats        *mq.Nats
	ServiceInfo micro.Info
}

// Init create the scaffold context
func (s *ScaffoldContext) Init(nats *mq.Nats, serviceInfo micro.Info) *ScaffoldContext {

	*s = ScaffoldContext{
		Nats:        nats,
		ServiceInfo: serviceInfo,
	}

	return s
}

// PublishError convenience function to publish errors
func (s *ScaffoldContext) PublishError(req micro.Request, err error) {

	errErr := req.Error("500", err.Error(), nil)

	if errErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(errErr))
	}

	zap.L().Error("Error Occurred", zap.Error(err))
}

// DecodeRequest convenience function for decoding structure data
func (s *ScaffoldContext) DecodeRequest(req micro.Request, result any) error {

	// Structs are just json so we need to unmarshal it
	unmarshallError := json.Unmarshal(req.Data(), &result)

	return unmarshallError
}

func ScaffoldContextHelper(context any, ctx *ScaffoldContext, handler func(any, *ScaffoldContext, micro.Request)) micro.Handler {
	return micro.HandlerFunc(func(req micro.Request) {
		handler(context, ctx, req)
	})
}
