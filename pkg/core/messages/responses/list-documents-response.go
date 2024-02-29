package responses

import "github.com/r4stl1n/micro-serv/pkg/core/messages/base"

type ListDocumentsResponse struct {
	Documents  []base.Document
	Collection string
}
