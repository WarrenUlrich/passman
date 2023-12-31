package passman

import "encoding/gob"

func init() {
	gob.Register(AddRequest{})
	gob.Register(AddResponse{})

	gob.Register(ListRequest{})
	gob.Register(ListResponse{})

	gob.Register(GetRequest{})
	gob.Register(GetResponse{})

	gob.Register(UpdateRequest{})
	gob.Register(UpdateResponse{})

	gob.Register(DeleteRequest{})
	gob.Register(DeleteResponse{})

	gob.Register(LockRequest{})
	gob.Register(LockResponse{})
}
