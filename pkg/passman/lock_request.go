package passman

type LockRequest struct {
	Password string
}

type LockResponse struct {
	Success bool
	Error   string
}
