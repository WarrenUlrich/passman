package passman

type AddRequest struct {
	Service  string
	Username string
	Password string
	Notes    string
}

type AddResponse struct {
	Success bool
}
