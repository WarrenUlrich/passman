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

	var err error
	db, err = sql.Open("sqlite3", path)
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

func handleAddRequest(request passman.AddRequest) error {
	fmt.Printf("Received add request: %+v\n", request)

	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := db.Exec(
		"INSERT INTO passwords (service, username, password, notes, expiry) VALUES (?, ?, ?, ?, ?)",
		request.Service,
		request.Username,
		request.Password,
		request.Notes,
		request.Expiry,
	)

	if err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn) {
	var request interface{}
	if err := gob.NewDecoder(conn).Decode(&request); err != nil {
		panic(err)
	}

	var err error
	switch request.(type) {
	case passman.AddRequest:
		err = handleAddRequest(request.(passman.AddRequest))
	default:
	}

	if err != nil {
		panic(err)
	}
	
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
