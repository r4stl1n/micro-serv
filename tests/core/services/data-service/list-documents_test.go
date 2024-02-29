package data_service

import (
	"github.com/r4stl1n/micro-serv/internal/core/data-service/service"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/mediators"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

func TestListDocuments(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Create the service
	serviceMan, serviceManError := new(service.DataService).Init(consts.TestPrefixConst.ToString(), "test")

	if serviceManError != nil {
		zap.L().Fatal("failed to create data service", zap.Error(serviceManError))
	}

	// Start the service
	runError := serviceMan.Run()

	if runError != nil {
		zap.L().Fatal("failed to run data service", zap.Error(runError))
	}

	// Create a new nats client from the environment variables
	natsClient, natsClientError := new(mq.Nats).Init()

	if natsClientError != nil {
		t.Skip("nats not configured for test")
	}

	// Connect to the nats server
	connectError := natsClient.Connect()

	if connectError != nil {
		zap.L().Error("failed to connect to nats server")
		t.Fatal()
	}

	// Create a new cache mediator
	dataMediator := new(mediators.DataServiceMediator).Init(natsClient, consts.VersionV1Const, consts.TestPrefixConst.ToString())

	// Send the add symbol request
	addDocumentResponse, addDocumentResponseError := dataMediator.AddDocument("testData", base.DummyMessage{
		Text: "testData",
	})

	if addDocumentResponseError != nil {
		zap.L().Error("failed to add document", zap.Error(addDocumentResponseError))
		t.Fatal()
	}

	// Send the list message
	listDocumentsResponse, listDocumentsResponseError := dataMediator.ListDocuments("testData")

	if listDocumentsResponseError != nil {
		zap.L().Error("failed to list documents", zap.Error(listDocumentsResponseError))
		t.Fatal()
	}

	if len(listDocumentsResponse.Documents) <= 0 {
		zap.L().Error("no documents returned")
		t.Fatal()
	}

	for _, doc := range listDocumentsResponse.Documents {
		var testMessage base.DummyMessage
		decodeError := doc.Decode(&testMessage)
		if decodeError != nil {
			zap.L().Error("failed to decode listed documents", zap.Error(decodeError))
			t.Fatal()
		}
	}

	// Send a removal message
	_, removeDocumentResponseError := dataMediator.RemoveDocument("testData", addDocumentResponse.ID)

	if removeDocumentResponseError != nil {
		zap.L().Error("failed to remove document", zap.Error(removeDocumentResponseError))
		t.Fatal()
	}

	// Shut it all down
	serviceStopError := serviceMan.Stop()

	if serviceStopError != nil {
		zap.L().Error("failed to stop the service", zap.Error(serviceStopError))
		t.Fatal()
	}

	natsClient.Disconnect()

}
