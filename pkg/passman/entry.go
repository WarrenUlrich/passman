package passman

import "time"

type Entry struct {
	Service  string
	Username string
	Password string
	Notes    string
	Expiry   time.Time
}
