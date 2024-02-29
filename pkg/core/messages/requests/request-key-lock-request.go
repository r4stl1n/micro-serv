package requests

import "time"

type RequestKeyLockRequest struct {
	Key        string
	Expiration time.Duration
}
