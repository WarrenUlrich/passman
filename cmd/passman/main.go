package main

import (
	"flag"
	"fmt"
	"time"
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
	clientConn *ClientConn
)

func makeRequest(request interface{}) error {
	return nil
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

	expiryTime := time.Time{}
	if *expiry != "" {
		var err error
		expiryTime, err = time.Parse(time.RFC3339, *expiry)
		if err != nil {
			return err
		}
	}

	_, err := clientConn.Add(remainingArgs[0], remainingArgs[1], *pwd, *notes, expiryTime)
	if err != nil {
		return err
	}

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
	clientConn, err = NewClientConn(socketPath)
	if err != nil {
		panic(err)
	}

	if err = runCommand(command, args); err != nil {
		panic(err)
	}
}
