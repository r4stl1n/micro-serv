package mediators

import (
	"fmt"
	"github.com/r4stl1n/micro-serv/pkg/core/consts"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/requests"
	"github.com/r4stl1n/micro-serv/pkg/core/messages/responses"
	"github.com/r4stl1n/micro-serv/pkg/core/mq"
	"go.mongodb.org/mongo-driver/bson"
)

type DataServiceMediator struct {
	natsClient *mq.Nats
	prefix     string
	version    consts.VersionConst
}

// Init creates a new data mediator for interacting with the data service
func (d *DataServiceMediator) Init(nats *mq.Nats, version consts.VersionConst, prefix string) *DataServiceMediator {
	*d = DataServiceMediator{
		natsClient: nats,
		prefix:     prefix,
		version:    version,
	}

	return d
}

// AddDocument attempts to store a document in the data service
func (d *DataServiceMediator) AddDocument(collection string, data any) (responses.AddDocumentResponse, error) {
	var addDocumentResponse responses.AddDocumentResponse

	// Send the add symbol request
	addDocumentResponseError := d.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", d.version, d.prefix, consts.CoreDataEndpointAddDocument),
		requests.AddDocumentRequest{
			Collection: collection,
			Data:       data,
		}, &addDocumentResponse)

	return addDocumentResponse, addDocumentResponseError
}

// RemoveDocument attempts to remove a document stored in the data service
func (d *DataServiceMediator) RemoveDocument(collection string, documentId string) (responses.RemoveDocumentResponse, error) {
	var removeDocumentResponse responses.RemoveDocumentResponse
	// Send a removal message
	removeDocumentResponseError := d.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", d.version, d.prefix, consts.CoreDataEndpointRemoveDocument),
		requests.RemoveDocumentRequest{
			Id:         documentId,
			Collection: collection,
		}, &removeDocumentResponse)

	return removeDocumentResponse, removeDocumentResponseError
}

// FilterDocuments attempts to search documents in the data service
func (d *DataServiceMediator) FilterDocuments(collection string, filters bson.D) (responses.FilterDocumentsResponse, error) {
	var filterDocumentsResponse responses.FilterDocumentsResponse
	// Send the filter message
	filterDocumentsResponseError := d.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", d.version, d.prefix, consts.CoreDataEndpointFilterDocuments),
		requests.FilterDocumentsRequest{
			Collection: collection,
			Filters:    filters,
		}, &filterDocumentsResponse)

	return filterDocumentsResponse, filterDocumentsResponseError
}

// ListDocuments retrieves all documents in the collection
func (d *DataServiceMediator) ListDocuments(collection string) (responses.ListDocumentsResponse, error) {
	var listDocumentsResponse responses.ListDocumentsResponse
	// Send the list message
	listDocumentsResponseError := d.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", d.version, d.prefix, consts.CoreDataEndpointListDocuments),
		requests.ListDocumentsRequest{
			Collection: collection,
		}, &listDocumentsResponse)

	return listDocumentsResponse, listDocumentsResponseError
}

// RetrieveDocument attempts to retrieve a document from the data store
func (d *DataServiceMediator) RetrieveDocument(collection string, documentId string, output any) error {
	var retrieveDocumentResponse responses.RetrieveDocumentResponse

	retrieveDocumentResponseError := d.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", d.version, d.prefix, consts.CoreDataEndpointRetrieveDocument),
		requests.RetrieveDocumentRequest{
			ID:         documentId,
			Collection: collection,
		}, &retrieveDocumentResponse)

	if retrieveDocumentResponseError != nil {
		return retrieveDocumentResponseError
	}

	return retrieveDocumentResponse.Document.Decode(output)

}

// ReplaceDocument attempts to replace a document
func (d *DataServiceMediator) ReplaceDocument(collection string, documentId string, data any) (responses.ReplaceDocumentResponse, error) {
	var replaceDocumentResponse responses.ReplaceDocumentResponse

	// Send the add symbol request
	replaceDocumentResponseError := d.natsClient.SendRequest(
		fmt.Sprintf("%s.%s%s", d.version, d.prefix, consts.CoreDataEndpointReplaceDocument),
		requests.ReplaceDocumentRequest{
			ID:         documentId,
			Collection: collection,
			Data:       data,
		}, &replaceDocumentResponse)

	return replaceDocumentResponse, replaceDocumentResponseError
}
