package passman

type ListRequest struct {
	Query string
}

type ListResponse struct {
	Entries []Entry
}
