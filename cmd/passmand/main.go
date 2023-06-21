package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/user"

	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/warrenulrich/passman/pkg/passman"
)

const (
	socketPath = "/tmp/passmand.sock"
)

var (
	db *sql.DB
)

func initializeDatabase(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	const createTableSQL = `
		CREATE TABLE IF NOT EXISTS passwords (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service TEXT NOT NULL,
			username TEXT NOT NULL,
			password TEXT,
			notes TEXT,
			expiry TIMESTAMP
		);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	return nil
}

func handleAddRequest(request passman.AddRequest) (*passman.AddResponse, error) {
	fmt.Printf("Received add request: %+v\n", request)
	return nil, nil
}

func handleConnection(conn net.Conn) {
	var request interface{}
	if err := gob.NewDecoder(conn).Decode(&request); err != nil {
		panic(err)
	}

	var response interface{}
	var err error
	switch request.(type) {
	case passman.AddRequest:
		response, err = handleAddRequest(request.(passman.AddRequest))
	default:
	}

	if err != nil {
		panic(err)
	}

	_ = response
	// if err := gob.NewEncoder(conn).Encode(response); err != nil {
	// 	panic(err)
	// }

	if err := conn.Close(); err != nil {
		panic(err)
	}
}

func listen() error {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	fmt.Println("Listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		fmt.Println("Accepted connection")
		go handleConnection(conn)
	}
}

func main() {
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	err = initializeDatabase(
		filepath.Join(usr.HomeDir, ".passman", "passman.db"),
	)

	if err != nil {
		panic(err)
	}

	if err := listen(); err != nil {
		panic(err)
	}
}
