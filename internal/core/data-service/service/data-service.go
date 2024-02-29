package service

import (
	"fmt"
	"github.com/r4stl1n/micro-serv/internal/core/data-service/context"
	"github.com/r4stl1n/micro-serv/internal/core/data-service/handlers"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/db"
	"github.com/r4stl1n/micro-serv/pkg/core/scaff"
)

type DataService struct {
	scaffold       *scaff.ServiceScaffold
	dataCtx        *context.DataContext
	endpointPrefix string
}

// Init create a new DataService
func (d *DataService) Init(endpointPrefix string, dbName string) (*DataService, error) {

	// Create a new scaffold instance
	scaffold, scaffoldingError := new(scaff.ServiceScaffold).Init("data-service", "1.0.0", "Data Service")

	if scaffoldingError != nil {
		return nil, scaffoldingError
	}

	// Create mongodb connection
	mongo, mongoError := new(db.Mongo).Init()

	if mongoError != nil {
		return nil, mongoError
	}

	// Connect to mongodb
	mongoConnectError := mongo.Connect()
	if mongoConnectError != nil {
		return nil, mongoConnectError
	}

	*d = DataService{
		dataCtx:        new(context.DataContext).Init(dbName, mongo),
		scaffold:       scaffold,
		endpointPrefix: endpointPrefix,
	}
	return d, nil
}

// StartUp create the handlers and perform any additional first time needs
func (d *DataService) addHandlers() error {

	addDocumentHandlerError := d.scaffold.AddHandler(d.dataCtx,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDataEndpointAddDocument), handlers.AddDocumentHandler)

	if addDocumentHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDataEndpointAddDocument, addDocumentHandlerError)
	}

	removeDocumentHandlerError := d.scaffold.AddHandler(d.dataCtx,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDataEndpointRemoveDocument), handlers.RemoveDocumentHandler)

	if removeDocumentHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDataEndpointRemoveDocument, addDocumentHandlerError)
	}

	listDocumentsHandlerError := d.scaffold.AddHandler(d.dataCtx,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDataEndpointListDocuments), handlers.ListDocumentsHandler)

	if listDocumentsHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDataEndpointListDocuments, addDocumentHandlerError)
	}

	retrieveDocumentHandlerError := d.scaffold.AddHandler(d.dataCtx,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDataEndpointRetrieveDocument), handlers.RetrieveDocumentHandler)

	if retrieveDocumentHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDataEndpointRetrieveDocument, retrieveDocumentHandlerError)
	}

	replaceDocumentHandlerError := d.scaffold.AddHandler(d.dataCtx,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDataEndpointReplaceDocument), handlers.ReplaceDocumentHandler)

	if replaceDocumentHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDataEndpointReplaceDocument, replaceDocumentHandlerError)
	}

	filterDocumentsHandlerError := d.scaffold.AddHandler(d.dataCtx,
		fmt.Sprintf("%s%s", d.endpointPrefix, consts.CoreDataEndpointFilterDocuments), handlers.FilterDocumentsHandler)

	if filterDocumentsHandlerError != nil {
		return fmt.Errorf("failed to add %s%s handler: %s",
			d.endpointPrefix, consts.CoreDataEndpointFilterDocuments, filterDocumentsHandlerError)
	}

	return nil
}

// Run the service
func (d *DataService) Run() error {

	startError := d.scaffold.Start()

	if startError != nil {
		return fmt.Errorf("failed to start the service: %s", startError)
	}

	return d.addHandlers()
}

// Stop the service
func (d *DataService) Stop() error {
	stopError := d.scaffold.Stop()

	if stopError != nil {
		return fmt.Errorf("failed to stop the service: %s", stopError)
	}

	return nil
}

// IsRunning returns if the service is running
func (d *DataService) IsRunning() bool {
	return d.scaffold.IsRunning()
}
