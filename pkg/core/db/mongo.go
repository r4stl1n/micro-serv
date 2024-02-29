package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"reflect"
	"time"
)

// Mongo creates and manages the db message queue connection
type Mongo struct {
	config *MongoConfig
	client *mongo.Client
	error  error
}

// Init create a new db manager with a given config
func (m *Mongo) Init() (*Mongo, error) {

	config, configError := new(MongoConfig).FromEnv()

	if configError != nil {
		return nil, configError
	}

	*m = Mongo{
		config: config,
	}

	return m, nil
}

// GetConfig get the db config
func (m *Mongo) GetConfig() MongoConfig {
	return *m.config
}

// Connect to the db server
func (m *Mongo) Connect() error {

	// Check if our nats configuration is set properly
	if !m.config.IsSet() {
		return fmt.Errorf("nats not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Create the db connection string
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?retryWrites=true&w=majority", m.config.User, m.config.Pass, m.config.Host)

	// Create a connection to the db database
	client, clientError := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if clientError != nil {
		return clientError
	}

	// Test the mongodb database connection
	pingError := client.Ping(ctx, readpref.Primary())

	if pingError != nil {
		return pingError
	}

	zap.L().Info("connected to db server", zap.String("host", m.config.Host))

	m.client = client

	return nil
}

// Disconnect from the db server
func (m *Mongo) Disconnect() error {
	disconnectError := m.client.Disconnect(context.Background())

	if disconnectError != nil {
		return disconnectError
	}

	zap.L().Info("disconnected from db server", zap.String("host", m.config.Host))

	return nil
}

// GetClient retrieves the raw db connection
func (m *Mongo) GetClient() *mongo.Client {
	return m.client
}

// AddRecord inserts a record into the db database
func (m *Mongo) AddRecord(database string, collection string, record any) (string, error) {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve all documents from the collection
	addResult, addError := m.client.Database(database).Collection(collection).InsertOne(ctx, &record)

	if addError != nil {
		return "", addError
	}

	zap.L().Info("added db record", zap.String("collection", collection),
		zap.String("id", addResult.InsertedID.(primitive.ObjectID).Hex()),
		zap.Int("size", int(reflect.TypeOf(record).Size())))

	return addResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

// ReplaceRecord update a record in the db database
func (m *Mongo) ReplaceRecord(database string, collection string, id string, record any) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert our string id to a primitive hex id
	hexId, hexIdError := primitive.ObjectIDFromHex(id)

	if hexIdError != nil {
		return hexIdError
	}

	// Replace the document
	replace := m.client.Database(database).Collection(collection).FindOneAndReplace(ctx,
		bson.D{{"_id", hexId}}, &record)

	if replace.Err() != nil {
		return replace.Err()
	}

	zap.L().Info("replaced db record", zap.String("collection", collection), zap.String("id", id),
		zap.Int("size", int(reflect.TypeOf(record).Size())))

	return nil
}

// RemoveRecord deletes a record into the db database
func (m *Mongo) RemoveRecord(database string, collection string, id string) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert our string id to a primitive hex id
	hexId, hexIdError := primitive.ObjectIDFromHex(id)

	if hexIdError != nil {
		return hexIdError
	}

	// Retrieve all documents from the collection
	deleteResult, deleteError := m.client.Database(database).Collection(collection).DeleteOne(ctx, bson.M{"_id": hexId})

	if deleteError != nil {
		return deleteError
	}

	if deleteResult.DeletedCount <= 0 {
		return fmt.Errorf("no records found to delete")
	}

	zap.L().Info("removed db record", zap.String("collection", collection), zap.String("id", id))

	return nil
}

// RetrieveAll get all documents in a collection
// Needs to be an array of interfaces
func (m *Mongo) RetrieveAll(database string, collection string, result any) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve all documents from the collection
	cursor, findError := m.client.Database(database).Collection(collection).Find(ctx, bson.D{})

	if findError != nil {
		return findError
	}

	// Process the retrieved to the result struct
	allError := cursor.All(ctx, result)

	if allError != nil {
		return allError
	}

	zap.L().Info("retrieved all db records", zap.String("collection", collection))

	return nil
}

// RetrieveFiltered get all documents in a collection
// Needs to be an array of interfaces
func (m *Mongo) RetrieveFiltered(database string, collection string, filter bson.D, result any) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filterBytes, filterMarshalError := bson.MarshalExtJSON(filter, false, false)

	if filterMarshalError != nil {
		return filterMarshalError
	}

	// Retrieve all documents from filter
	cursor, findError := m.client.Database(database).Collection(collection).Find(ctx, filter)

	if findError != nil {
		return findError
	}

	// Process the retrieved to the result struct
	allError := cursor.All(ctx, result)

	if allError != nil {
		return allError
	}

	zap.L().Info("retrieved filtered db records", zap.String("collection", collection),
		zap.String("filter", string(filterBytes[:])))

	return nil
}

// RetrieveOne get a single document in a collection
func (m *Mongo) RetrieveOne(database string, collection string, id string, result any) error {

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert our string id to a primitive hex id
	hexId, hexIdError := primitive.ObjectIDFromHex(id)

	if hexIdError != nil {
		return hexIdError
	}

	// Retrieve all documents from filter
	cursor := m.client.Database(database).Collection(collection).FindOne(ctx, bson.M{"_id": hexId})

	if cursor.Err() != nil {
		return cursor.Err()
	}

	resultError := cursor.Decode(result)

	if resultError != nil {
		return resultError
	}

	zap.L().Info("retrieved one db record", zap.String("collection", collection), zap.String("id", id))

	return nil
}
