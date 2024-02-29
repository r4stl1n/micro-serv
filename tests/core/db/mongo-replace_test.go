package services

import (
	"github.com/r4stl1n/micro-serv/pkg/core/db"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"github.com/r4stl1n/micro-serv/pkg/core/util"
	"go.uber.org/zap"
	"testing"
)

func TestMongoReplace(t *testing.T) {

	// Create a new logger
	new(util.Logger).Init()

	// Being lazy and using this to check if nats is running
	mongoClient, mongoClientError := new(db.Mongo).Init()

	if mongoClientError != nil {
		t.Skip("mongodb not configured for test")
	}

	connectError := mongoClient.Connect()

	if connectError != nil {
		zap.L().Error("failed to connect to mongodb", zap.Error(connectError))
		t.Fatal()
	}

	// Insert a new record into the db database
	recordId, insertError := mongoClient.AddRecord("test", "test", base.DummyMessage{Text: "herpderp"})
	if insertError != nil {
		zap.L().Error("failed to insert record to mongodb", zap.Error(insertError))
		t.Fatal()
	}

	zap.L().Info("record created", zap.String("id", recordId))

	// Replace the freshly created record
	replaceError := mongoClient.ReplaceRecord("test", "test", recordId, base.DummyMessage{Text: "merpherp"})

	if replaceError != nil {
		zap.L().Error("failed to replace single record in mongodb", zap.Error(replaceError))
		t.Fatal()
	}

	zap.L().Info("record replaced", zap.String("id", recordId), zap.String("text", "merpherp"))

	// Retrieve the updated record from the database
	var resultConfirm base.DummyMessage
	singleRecordConfirmError := mongoClient.RetrieveOne("test", "test", recordId, &resultConfirm)

	if singleRecordConfirmError != nil {
		zap.L().Error("failed to retrieve single record from mongodb", zap.Error(singleRecordConfirmError))
		t.Fatal()
	}

	// Verify the text on the record has been changed
	if resultConfirm.Text != "merpherp" {
		zap.L().Error("failed to update record, text remains the same")
		t.Fatal()
	}

	// Delete the record from the database
	deleteError := mongoClient.RemoveRecord("test", "test", recordId)
	if deleteError != nil {
		zap.L().Error("failed to delete record in mongodb", zap.Error(deleteError))
		t.Fatal()
	}

	zap.L().Info("record deleted", zap.String("id", recordId))

	// Disconnect from the db database
	disconnectError := mongoClient.Disconnect()

	if disconnectError != nil {
		zap.L().Error("failed to disconnect from mongodb", zap.Error(disconnectError))
		t.Fatal()
	}
}
