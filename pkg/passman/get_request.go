package passman

type GetRequest struct {
	Service  string
	Username string
}

type GetResponse struct {
	Password string
	Notes    string
	Expiry   string
}
