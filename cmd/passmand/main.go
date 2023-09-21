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

func handleAddRequest(conn net.Conn, request passman.AddRequest) error {
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

	resp := passman.AddResponse{
		Success: true,
	}

	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
		return err
	}

	return nil
}

func handleListRequest(conn net.Conn, request passman.ListRequest) error {
	fmt.Printf("Received list request: %+v\n", request)

	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	rows, err := db.Query(
		"SELECT service, username, password, notes, expiry FROM passwords",
	)

	if err != nil {
		return err
	}

	var entries []passman.Entry
	for rows.Next() {
		var entry passman.Entry
		if err := rows.Scan(
			&entry.Service,
			&entry.Username,
			&entry.Password,
			&entry.Notes,
			&entry.Expiry,
		); err != nil {
			return err
		}

		entries = append(entries, entry)
	}

	resp := passman.ListResponse{
		Entries: entries,
	}

	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
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
		err = handleAddRequest(conn, request.(passman.AddRequest))
	case passman.ListRequest:
		err = handleListRequest(conn, request.(passman.ListRequest))
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

	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		dataDir = filepath.Join(usr.HomeDir, ".local", "share")
	}

	err = initializeDatabase(
		filepath.Join(dataDir, "passman", "passman.db"),
	)

	if err != nil {
		panic(err)
	}

	if err := listen(); err != nil {
		panic(err)
	}
}
