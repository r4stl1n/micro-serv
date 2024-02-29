package responses

import (
	"encoding/json"
)

type RetrieveCacheResponse struct {
	Key  string
	Data any
}

func (r RetrieveCacheResponse) DecodeData(result any) error {

	// convert map to json
	marshalledBytes, marshalError := json.Marshal(r.Data)

	if marshalError != nil {
		return marshalError
	}

	return json.Unmarshal(marshalledBytes, result)
}
