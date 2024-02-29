package context

import (
	"github.com/r4stl1n/micro-serv/pkg/core/db"
)

type DataContext struct {
	DbName string
	Mongo  *db.Mongo
}

// Init create the scaffold context
func (s *DataContext) Init(dbName string, mongo *db.Mongo) *DataContext {

	*s = DataContext{
		DbName: dbName,
		Mongo:  mongo,
	}

	return s
}
