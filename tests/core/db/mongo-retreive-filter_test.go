package services

import (
	"github.com/r4stl1n/micro-serv/pkg/core/db"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"testing"
)

func TestMongoRetrieveFilter(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Being lazy and using this to check if nats is running
	mongoClient, mongoClientError := new(db.Mongo).Init()

	if mongoClientError != nil {
		t.Skip("mongodb not configured for test")
	}

	// Connect to the mongodb database
	connectError := mongoClient.Connect()

	if connectError != nil {
		zap.L().Error("failed to connect to mongodb", zap.Error(connectError))
		t.Fatal()
	}

	// Insert one record into the mongodb database
	recordId, insertError := mongoClient.AddRecord("test", "test", base.DummyMessage{Text: "herpderp"})
	if insertError != nil {
		zap.L().Error("failed to insert record to mongodb", zap.Error(insertError))
		t.Fatal()
	}

	zap.L().Info("record created 1", zap.String("id", recordId))

	// Insert another record into the mongodb database
	recordId2, insertError2 := mongoClient.AddRecord("test", "test", base.DummyMessage{Text: "merpherp"})
	if insertError2 != nil {
		zap.L().Error("failed to insert second record to mongodb", zap.Error(insertError2))
		t.Fatal()
	}

	zap.L().Info("record created 2", zap.String("id", recordId2))

	// Use a filter and retrieve text messages with only herpderp in the text
	var results []base.DummyMessage
	filteredRecordsError := mongoClient.RetrieveFiltered("test", "test", bson.D{{"text", "herpderp"}}, &results)

	if filteredRecordsError != nil {
		zap.L().Error("failed to retrieve single record from mongodb", zap.Error(filteredRecordsError))
		t.Fatal()
	}

	// Record the amount of records saved
	zap.L().Info("filtered records retrieved", zap.Int("amount", len(results)))

	// Delete the first record from the db database
	deleteError := mongoClient.RemoveRecord("test", "test", recordId)
	if deleteError != nil {
		zap.L().Error("failed to delete record in mongodb", zap.Error(deleteError))
		t.Fatal()
	}

	zap.L().Info("record 1 deleted", zap.String("id", recordId))

	// Delete the second record from the db database
	deleteError2 := mongoClient.RemoveRecord("test", "test", recordId2)
	if deleteError2 != nil {
		zap.L().Error("failed to delete record in mongodb", zap.Error(deleteError2))
		t.Fatal()
	}

	zap.L().Info("record 2 deleted", zap.String("id", recordId))

	// Disconnect from the mongodb database
	disconnectError := mongoClient.Disconnect()

	if disconnectError != nil {
		zap.L().Error("failed to disconnect from mongodb", zap.Error(disconnectError))
		t.Fatal()
	}
}
