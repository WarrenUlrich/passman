package passman

import "encoding/gob"

func init() {
	gob.Register(AddRequest{})
	gob.Register(AddResponse{})
	gob.Register(ListRequest{})
	gob.Register(ListResponse{})
}
