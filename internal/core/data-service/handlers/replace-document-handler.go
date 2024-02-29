package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/data-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func ReplaceDocumentHandler(dataContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	dContext := dataContext.(*context.DataContext)

	// Extract request message
	var replaceDocumentRequest requests.ReplaceDocumentRequest

	decodeError := sContext.DecodeRequest(microReq, &replaceDocumentRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Replace the record requested
	removeError := dContext.Mongo.ReplaceRecord(dContext.DbName, replaceDocumentRequest.Collection,
		replaceDocumentRequest.ID, &replaceDocumentRequest.Data)

	if removeError != nil {
		sContext.PublishError(microReq, removeError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.ReplaceDocumentResponse{ID: replaceDocumentRequest.ID})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}
}
