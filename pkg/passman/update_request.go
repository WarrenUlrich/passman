package passman

type UpdateRequest struct {
	Service  string
	Username string
	Password string
	Notes    string
}

type UpdateResponse struct {
	Success bool
	Error   string
}
