package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/warrenulrich/passman/pkg/passman"
)

const (
	helpMessage = `usage: passman [options] <command>

Options:
    -h, --help          Show help message and exit.
    --version           Show program's version number and exit.

Commands:
    add <service> <username> Add a new entry to the vault.
        -p, --password       Specify the password.
        -n, --notes          Add notes to the entry.
        -e, --expiry         Set an expiry date for the password.

    update <service> <username> Update an existing entry.
        Options are the same as for the 'add' command.

    get <service> <username> Retrieve an entry.
        -c, --clip           Copy the password to the clipboard.
        -s, --show           Display the password on the terminal.

    list [query]             List all entries or those that match a query.

    delete <service> <username> Delete an entry from the vault.

    generate                 Generate a secure password.
        -l, --length         Specify the length of the password.
        -s, --symbols        Include symbols.
        -n, --numbers        Include numbers.
        -u, --uppercase      Include uppercase letters.
        -x, --exclude        Exclude specific characters.

    export <file>            Export the vault data to a file.
        -f, --format         Specify the export format (CSV, JSON).

    import <file>            Import data from a file to the vault.
        Options are the same as for the 'export' command.

    change-master            Change the master password of the vault.

    config                   Configure settings like auto-lock time, default password length, etc.

    lock                     Manually lock the password vault.

    health                   Check the health of passwords in the vault.

    search <query>           Search the vault for the given query.`

	versionMessage = "passman version 1.0.0"

	socketPath = "/tmp/passmand.sock"
)

var (
	clientConn net.Conn
)

func writeRequest(request interface{}) error {
	return gob.NewEncoder(clientConn).Encode(&request)
}

func readResponse() (interface{}, error) {
	return nil, nil
}

func addCommand(args []string) error {
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	pwd := flags.String("p", "", "Specify the password.")
	notes := flags.String("n", "", "Add notes to the entry.")
	expiry := flags.String("e", "", "Set an expiry date for the password.")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remainingArgs := flags.Args()
	if len(remainingArgs) != 2 {
		for _, arg := range remainingArgs {
			fmt.Println(arg)
		}

		return fmt.Errorf("invalid number of arguments, need 2, got %d", len(remainingArgs))
	}

	expiryTime := time.Now()
	if *expiry != "" {
		var err error
		expiryTime, err = time.Parse(time.RFC3339, *expiry)
		if err != nil {
			return err
		}
	}

	request := passman.AddRequest{
		Service:  remainingArgs[0],
		Username: remainingArgs[1],
		Password: *pwd,
		Notes:    *notes,
		Expiry:   expiryTime,
	}

	_ = request
	if err := writeRequest(request); err != nil {
		return err
	}

	// if resp, err := readResponse(); err != nil {
	// 	fmt.Printf("Response: %v", resp)
	// }

	return nil
}

func runCommand(cmd string, args []string) error {
	switch cmd {
	case "add":
		return addCommand(args)
	}

	return fmt.Errorf("unknown command: %s", cmd)
}

func main() {
	flag.Usage = func() {
		fmt.Println(helpMessage)
	}

	help := flag.Bool("h", false, "")
	version := flag.Bool("version", false, "")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *version {
		fmt.Println(versionMessage)
		return
	}

	remainingArgs := flag.Args()

	if len(remainingArgs) == 0 {
		flag.Usage()
		return
	}

	command := remainingArgs[0]
	args := remainingArgs[1:]

	var err error
	clientConn, err = net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}

	if err = runCommand(command, args); err != nil {
		panic(err)
	}
}
