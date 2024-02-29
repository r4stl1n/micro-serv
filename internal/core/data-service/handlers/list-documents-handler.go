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

func ListDocumentsHandler(dataContext any, sContext *scaff.ScaffoldContext, microReq micro.Request) {

	// Cast the context for use
	dContext := dataContext.(*context.DataContext)

	// Extract request message
	var listDocumentsRequest requests.ListDocumentsRequest

	decodeError := sContext.DecodeRequest(microReq, &listDocumentsRequest)

	if decodeError != nil {
		sContext.PublishError(microReq, decodeError)
		return
	}

	// Retrieve all documents in the collection
	var documents []primitive.D
	retrieveError := dContext.Mongo.RetrieveAll(dContext.DbName, listDocumentsRequest.Collection, &documents)

	if retrieveError != nil {
		sContext.PublishError(microReq, retrieveError)
		return
	}

	var respDocuments []base.Document

	for _, doc := range documents {

		preparedDoc, prepareDocError := new(base.Document).Init().BsonPrepare(doc)

		if prepareDocError != nil {
			zap.L().Error("failed to prepare document", zap.Error(prepareDocError))
		}

		respDocuments = append(respDocuments, preparedDoc)
	}

	// Respond with a response message
	respErr := microReq.RespondJSON(&responses.ListDocumentsResponse{Documents: respDocuments,
		Collection: listDocumentsRequest.Collection})

	if respErr != nil {
		zap.L().Error("error publishing nats error message", zap.Error(respErr))
		return
	}

}
