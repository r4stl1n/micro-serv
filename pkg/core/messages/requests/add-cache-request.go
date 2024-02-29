package requests

import "time"

type AddCacheRequest struct {
	Key        string
	Data       any
	Expiration time.Duration
}
