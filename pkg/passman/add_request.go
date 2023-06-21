package passman

import "time"

type AddRequest struct {
	Service  string
	Username string
	Password string
	Notes    string
	Expiry   time.Time
}

type AddResponse struct {
	Success bool
	Message string
}
