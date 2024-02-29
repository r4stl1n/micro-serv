package handlers

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/r4stl1n/micro-serv/internal/core/data-service/context"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func RetrieveDocumentHandler(dataContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	dContext := dataContext.(*context.DataContext)

	// Extract request message
	var retrieveDocumentRequest requests.RetrieveDocumentRequest

	decodeError := sContext.DecodeRequest(microReq, &retrieveDocumentRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	var doc primitive.D

	// Retrieve the record requested
	removeError := dContext.Mongo.RetrieveOne(dContext.DbName, retrieveDocumentRequest.Collection,
		retrieveDocumentRequest.ID, &doc)

	if removeError != nil {
		sContext.PublishError(microReq, removeError)
		return
	}

	preparedDoc, prepareDocError := new(base.Document).Init().BsonPrepare(doc)

	if prepareDocError != nil {
		sContext.PublishError(microReq, prepareDocError)
		return
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.RetrieveDocumentResponse{ID: retrieveDocumentRequest.ID,
		Collection: retrieveDocumentRequest.Collection, Document: preparedDoc})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}
}
