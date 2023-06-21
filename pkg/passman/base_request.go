package passman

type RequestType int

type ResponseType int

const (
	AddRequestType RequestType = iota
)

const (
	AddResponseType ResponseType = iota
)

type BaseRequest struct {
	Type RequestType
}


type BaseResponse struct {
	Type ResponseType
	
}

