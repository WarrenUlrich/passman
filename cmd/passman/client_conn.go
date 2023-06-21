package main

import (
	"encoding/gob"
	"net"
	"time"

	"github.com/warrenulrich/passman/pkg/passman"
)

type ClientConn struct {
	conn net.Conn
}

func NewClientConn(path string) (*ClientConn, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &ClientConn{
		conn: conn,
	}, nil
}

func (c *ClientConn) Close() error {
	return c.conn.Close()
}

func (c *ClientConn) sendRequest(request interface{}) error {
	return gob.NewEncoder(c.conn).Encode(&request)
}

func (c *ClientConn) readResponse() (interface{}, error) {
	var response interface{}
	if err := gob.NewDecoder(c.conn).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientConn) Add(service, username, password, notes string, expiry time.Time) (*passman.AddResponse, error) {
	addRequest := passman.AddRequest{
		Service:  service,
		Username: username,
		Password: password,
		Notes:    notes,
		Expiry:   expiry,
	}

	if err := c.sendRequest(addRequest); err != nil {
		return nil, err
	}

	var addResponse passman.AddResponse
	if err := gob.NewDecoder(c.conn).Decode(&addResponse); err != nil {
		return nil, err
	}

	return &addResponse, nil
}
