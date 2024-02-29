package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/data-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.uber.org/zap"
)

func AddDocumentHandler(dataContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	dContext := dataContext.(*context.DataContext)

	// Extract request message
	var addDocumentRequest requests.AddDocumentRequest

	decodeError := sContext.DecodeRequest(microReq, &addDocumentRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Add the data to mongo database
	id, addError := dContext.Mongo.AddRecord(dContext.DbName, addDocumentRequest.Collection, addDocumentRequest.Data)

	if addError != nil {
		sContext.PublishError(microReq, addError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.AddDocumentResponse{ID: id,
		Collection: addDocumentRequest.Collection})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
