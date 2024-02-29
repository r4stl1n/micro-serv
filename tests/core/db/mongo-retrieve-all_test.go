package services

import (
	"github.com/r4stl1n/micro-serv/pkg/core/db"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

func TestMongoRetrieveAll(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Being lazy and using this to check if nats is running
	mongoClient, mongoClientError := new(db.Mongo).Init()

	if mongoClientError != nil {
		t.Skip("mongodb not configured for test")
	}

	// Connect to the db database
	connectError := mongoClient.Connect()

	if connectError != nil {
		zap.L().Error("failed to connect to mongodb", zap.Error(connectError))
		t.Fatal()
	}

	// Insert a record into the db database
	recordId, insertError := mongoClient.AddRecord("test", "test", base.DummyMessage{Text: "herpderp"})
	if insertError != nil {
		zap.L().Error("failed to insert record to mongodb", zap.Error(insertError))
		t.Fatal()
	}

	zap.L().Info("record created", zap.String("id", recordId))

	// Retrieve all text messages in the collection
	var results []base.DummyMessage

	allRecordsError := mongoClient.RetrieveAll("test", "test", &results)

	if allRecordsError != nil {
		zap.L().Error("failed to retrieve single record from mongodb", zap.Error(allRecordsError))
		t.Fatal()
	}

	zap.L().Info("all records retrieved", zap.Int("amount", len(results)))

	// Delete the record from the db database
	deleteError := mongoClient.RemoveRecord("test", "test", recordId)
	if deleteError != nil {
		zap.L().Error("failed to delete record in mongodb", zap.Error(deleteError))
		t.Fatal()
	}

	zap.L().Info("record deleted", zap.String("id", recordId))

	// Disconnect from the db client
	disconnectError := mongoClient.Disconnect()

	if disconnectError != nil {
		zap.L().Error("failed to disconnect from mongodb", zap.Error(disconnectError))
		t.Fatal()
	}
}
