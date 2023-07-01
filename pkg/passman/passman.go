package passman

import "encoding/gob"

type BaseResponse struct {
	Error error
}


func init() {
	gob.Register(BaseResponse{})
	gob.Register(AddRequest{})
	gob.Register(AddResponse{})
}
