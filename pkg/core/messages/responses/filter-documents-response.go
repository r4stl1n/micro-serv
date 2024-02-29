package responses

import (
	"github.com/r4stl1n/micro-serv/pkg/core/messages/base"
	"go.mongodb.org/mongo-driver/bson"
)

type FilterDocumentsResponse struct {
	Filter     bson.D
	Documents  []base.Document
	Collection string
}
