package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func EchoHandler(_ any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Extract the request message
	var textMessage base.DummyMessage

	decodeError := sContext.DecodeRequest(microReq, &textMessage)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Respond with the same text message
	respErr := microReq.RespondJSON(&textMessage)

	if respErr != nil {
		zap.L().Error("error publishing nats message", zap.Error(respErr))
		return
	}

	zap.L().Info("received message", zap.String("msg", textMessage.Text))

}
