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
			UNIQUE(service, username)
		);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	return nil
}

func handleAddRequest(conn net.Conn, req passman.AddRequest) error {
	_, err := db.Exec(
		"INSERT INTO passwords (service, username, password, notes) VALUES (?, ?, ?, ?)",
		req.Service,
		req.Username,
		req.Password,
		req.Notes,
	)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: passwords.service, passwords.username" {
			fmt.Println("error:", err.Error())

			resp := passman.AddResponse{
				Success: false,
				Error:   "Entry already exists",
			}

			if err := gob.NewEncoder(conn).Encode(resp); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	resp := passman.AddResponse{
		Success: true,
	}

	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
		return err
	}

	return nil
}

func handleListRequest(conn net.Conn, req passman.ListRequest) error {
	rows, err := db.Query(
		"SELECT service, username, password, notes FROM passwords",
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

func handleGetRequest(conn net.Conn, req passman.GetRequest) error {
	var entry passman.Entry
	if err := db.QueryRow(
		"SELECT service, username, password, notes FROM passwords WHERE service = ? AND username = ?",
		req.Service,
		req.Username,
	).Scan(
		&entry.Service,
		&entry.Username,
		&entry.Password,
		&entry.Notes,
	); err != nil {
		return err
	}

	resp := passman.GetResponse{
		Password: entry.Password,
		Notes:    entry.Notes,
	}

	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
		return err
	}

	return nil
}

func handleUpdateRequest(conn net.Conn, req passman.UpdateRequest) error {
	_, err := db.Exec(
		"UPDATE passwords SET password = ?, notes = ? WHERE service = ? AND username = ?",
		req.Password,
		req.Notes,
		req.Service,
		req.Username,
	)

	if err != nil {
		return err
	}

	resp := passman.UpdateResponse{
		Success: true,
	}

	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
		return err
	}

	return nil
}

func handleDeleteRequest(conn net.Conn, req passman.DeleteRequest) error {
	_, err := db.Exec(
		"DELETE FROM passwords WHERE service = ? AND username = ?",
		req.Service,
		req.Username,
	)

	if err != nil {
		return err
	}

	resp := passman.DeleteResponse{
		Success: true,
	}

	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
		return err
	}

	return nil
}

func handleLockRequest(conn net.Conn, req passman.LockRequest) error {
	// TODO: Implement encrypting the database with the master password

	var resp passman.LockResponse
	if err := gob.NewEncoder(conn).Encode(resp); err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn) {
	var req interface{}
	if err := gob.NewDecoder(conn).Decode(&req); err != nil {
		if err.Error() == "EOF" {
			return
		}

		panic(err)
	}

	var err error
	switch r := req.(type) {
	case passman.AddRequest:
		err = handleAddRequest(conn, r)
	case passman.ListRequest:
		err = handleListRequest(conn, r)
	case passman.GetRequest:
		err = handleGetRequest(conn, r)
	case passman.UpdateRequest:
		err = handleUpdateRequest(conn, r)
	case passman.DeleteRequest:
		err = handleDeleteRequest(conn, r)
	case passman.LockRequest:
		err = handleLockRequest(conn, r)
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
