package passman

type DeleteRequest struct {
	Service  string
	Username string
}

type DeleteResponse struct {
	Success bool
	Error   string
}
