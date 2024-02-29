package base

import (
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	ID  string `json:"id" bson:"_id,omitempty"`
	Raw string
}

// Init create the scaffold context
func (d *Document) Init() *Document {

	*d = Document{}

	return d
}

// BsonPrepare prepare the doc object with bson
func (d *Document) BsonPrepare(data any) (Document, error) {

	// Convert the bson data to a json string
	marshalBytes, marshalError := bson.Marshal(data)

	if marshalError != nil {
		return Document{}, marshalError
	}

	// Store the json string
	d.Raw = base64.StdEncoding.EncodeToString(marshalBytes[:])
	d.ID = data.(primitive.D)[0].Value.(primitive.ObjectID).Hex()

	return *d, marshalError
}

func (d *Document) Decode(result any) error {

	var temp map[string]any

	// Structs are just json so we need to unmarshal it
	// First we have to decode the bson data
	b64Decoded, b64DecodeError := base64.StdEncoding.DecodeString(d.Raw)

	if b64DecodeError != nil {
		return b64DecodeError
	}

	// Unmarshall the bson into a temp structure
	bsonUnmarshallError := bson.Unmarshal(b64Decoded, &temp)

	if bsonUnmarshallError != nil {
		return bsonUnmarshallError
	}

	// Final conversion to get the temp interface to json structure
	jsonString, jsonMarshallError := json.Marshal(&temp)

	if jsonMarshallError != nil {
		return jsonMarshallError
	}

	return json.Unmarshal(jsonString, &result)
}
