package requests

import "go.mongodb.org/mongo-driver/bson"

type FilterDocumentsRequest struct {
	Collection string
	Filters    bson.D
}
