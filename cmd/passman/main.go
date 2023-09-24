package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

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

    config                   Configure settings.

	lock                     Manually lock the password vault.

	unlock					 Manually unlock the password vault.`

	versionMessage = "passman version 0.0.1"

	socketPath = "/tmp/passmand.sock"
)

func getClient() (net.Conn, error) {
	return net.Dial("unix", socketPath)
}

func writeRequest(conn net.Conn, req interface{}) error {
	return gob.NewEncoder(conn).Encode(&req)
}

func readResponse[T any](conn net.Conn) (*T, error) {
	var resp T
	if err := gob.NewDecoder(conn).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func generatePassword(length int, symbols, numbers, uppcase bool, exclude []rune) string {
	const lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	const uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const symbolsSet = "!@#$%^&*()-_=+{}[]|:;<>,.?/"
	const numbersSet = "0123456789"

	// Start with lowercase letters.
	var charset string = lowercaseLetters

	// Add other character sets based on the flags.
	if uppcase {
		charset += uppercaseLetters
	}
	if symbols {
		charset += symbolsSet
	}
	if numbers {
		charset += numbersSet
	}

	// Remove excluded characters.
	for _, ex := range exclude {
		charset = strings.ReplaceAll(charset, string(ex), "")
	}

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(result)
}

func addCommand(args []string) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	pwd := flags.String("p", "", "Specify the password.")
	notes := flags.String("n", "", "Add notes to the entry.")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remainingArgs := flags.Args()

	if len(remainingArgs) < 2 {
		return errors.New("insufficient non-flag arguments")
	}

	service := remainingArgs[0]
	username := remainingArgs[1]

	request := passman.AddRequest{
		Service:  service,
		Username: username,
		Password: *pwd,
		Notes:    *notes,
	}

	fmt.Println("Request:", request)

	client, err := getClient()
	if err != nil {
		return err
	}

	if err = gob.NewEncoder(client).Encode(request); err != nil {
		return err
	}

	var response passman.AddResponse
	if err = gob.NewDecoder(client).Decode(&response); err != nil {
		return err
	}

	fmt.Printf("Response: %+v\n", response)
	return nil
}

func generateCommand(args []string) error {
	flags := flag.NewFlagSet("generate", flag.ExitOnError)

	length := flags.Int("l", 10, "")
	numbers := flags.Bool("n", false, "")
	symbols := flags.Bool("s", false, "")
	uppercase := flags.Bool("u", false, "")
	clip := flags.Bool("c", false, "")

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	pass := generatePassword(*length, *symbols, *numbers, *uppercase, []rune{})
	if *clip {
		// TODO: save pass to clipboard
	}

	fmt.Println(pass)
	return nil
}

func listCommand(args []string) error {
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	request := passman.ListRequest{
		Query: query,
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	// var i interface{} = request
	// if err = gob.NewEncoder(client).Encode(&i); err != nil {
	// 	return err
	// }

	if err = writeRequest(client, request); err != nil {
		return err
	}

	resp, err := readResponse[passman.ListResponse](client)
	if err != nil {
		return err
	}

	for _, entry := range resp.Entries {
		fmt.Printf("%s\t%s\t%s\t%s\n", entry.Service, entry.Username, entry.Password, entry.Notes)
	}

	return nil
}

func getCommand(args []string) error {
	flags := flag.NewFlagSet("get", flag.ExitOnError)

	clip := flags.Bool("c", false, "")
	show := flags.Bool("s", false, "")

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	if len(flags.Args()) < 2 {
		return errors.New("insufficient non-flag arguments")
	}

	service := flags.Args()[0]
	username := flags.Args()[1]

	request := passman.GetRequest{
		Service:  service,
		Username: username,
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if err = writeRequest(client, request); err != nil {
		return err
	}

	resp, err := readResponse[passman.GetResponse](client)
	if err != nil {
		return err
	}

	if *clip {
		// TODO: save pass to clipboard
	}

	if *show {
		fmt.Println("Pass:", resp.Password)
	}

	return nil
}

func updateCommand(args []string) error {
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	pwd := flags.String("p", "", "Specify the password.")
	notes := flags.String("n", "", "Add notes to the entry.")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remainingArgs := flags.Args()

	if len(remainingArgs) < 2 {
		return errors.New("insufficient non-flag arguments")
	}

	service := remainingArgs[0]
	username := remainingArgs[1]

	request := passman.UpdateRequest{
		Service:  service,
		Username: username,
		Password: *pwd,
		Notes:    *notes,
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if err = writeRequest(client, request); err != nil {
		return err
	}

	resp, err := readResponse[passman.UpdateResponse](client)
	if err != nil {
		return err
	}

	fmt.Printf("Response: %+v\n", resp)
	return nil
}

func deleteCommand(args []string) error {
	req := passman.DeleteRequest{
		Service:  args[0],
		Username: args[1],
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if err = writeRequest(client, req); err != nil {
		return err
	}

	resp, err := readResponse[passman.DeleteResponse](client)
	if err != nil {
		return err
	}

	fmt.Printf("Response: %+v\n", resp)
	return nil
}

func lockCommand(args []string) error {
	fmt.Println("Enter master password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	req := passman.LockRequest{
		Password: string(password),
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if err = writeRequest(client, req); err != nil {
		return err
	}

	resp, err := readResponse[passman.LockResponse](client)
	if err != nil {
		return err
	}

	fmt.Printf("Response: %+v\n", resp)
	return nil
}

func runCommand(cmd string, args []string) error {
	switch cmd {
	case "add":
		return addCommand(args)
	case "list":
		return listCommand(args)
	case "generate":
		return generateCommand(args)
	case "get":
		return getCommand(args)
	case "update":
		return updateCommand(args)
	case "delete":
		return deleteCommand(args)
	case "lock":
		return lockCommand(args)
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

	// var err error
	// clientConn, err = net.Dial("unix", socketPath)
	// if err != nil {
	// 	panic(err)
	// }

	if err := runCommand(command, args); err != nil {
		panic(err)
	}
}
