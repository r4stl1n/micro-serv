package responses

import "github.com/r4stl1n/micro-serv/pkg/core/messages/base"

type RetrieveDocumentResponse struct {
	ID         string
	Collection string
	Document   base.Document
}
