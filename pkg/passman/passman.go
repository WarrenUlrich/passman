package passman

import "encoding/gob"

func init() {
	gob.Register(AddRequest{})
	gob.Register(AddResponse{})
}
