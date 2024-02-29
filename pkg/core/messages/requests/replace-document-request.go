package requests

type ReplaceDocumentRequest struct {
	ID         string
	Data       any
	Collection string
}
