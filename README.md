# Passman - A Command-Line Password Manager in Go

Passman is a secure and easy-to-use command-line password manager written in Go.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Adding an Entry](#adding-an-entry)
  - [Updating an Entry](#updating-an-entry)
  - [Retrieving an Entry](#retrieving-an-entry)
  - [Listing Entries](#listing-entries)
  - [Deleting an Entry](#deleting-an-entry)
  - [Generating Passwords](#generating-passwords)
  - [Exporting and Importing Data](#exporting-and-importing-data)
  - [Changing the Master Password](#changing-the-master-password)
  - [Configuring Passman](#configuring-passman)
  - [Locking and Unlocking the Vault](#locking-and-unlocking-the-vault)
- [Options](#options)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Secure Storage**: Passman securely stores your passwords and sensitive information using strong encryption.

- **Easy-to-Use**: The command-line interface is user-friendly and designed for ease of use.

- **Password Generation**: Passman can generate strong and random passwords for you.

- **Flexible Import/Export**: Import and export your vault data in various formats like CSV and JSON.

- **Lock and Unlock**: Manually lock and unlock your password vault for enhanced security.

## Usage

Passman offers a wide range of commands and options to help you manage your passwords efficiently. Here are some common use cases:

### Adding an Entry

```sh
passman add <service> <username> -p <password> -n <notes>
```

Use the `add` command to add a new entry to the vault. You can specify the password and add optional notes.

### Updating an Entry

```sh
passman update <service> <username> -p <password> -n <notes>
```

Update an existing entry in the vault with the `update` command. Options are the same as for the 'add' command.

### Retrieving an Entry

```sh
passman get <service> <username> -c -s
```

Retrieve an entry using the `get` command. Use the `-c` option to copy the password to the clipboard and `-s` to display the password on the terminal.

### Listing Entries

```sh
passman list [query]
```

List all entries or those that match a query with the `list` command.

### Deleting an Entry

```sh
passman delete <service> <username>
```

Delete an entry from the vault using the `delete` command.

### Generating Passwords

```sh
passman generate -l <length> -s -n -u -x <exclude>
```

Generate secure passwords with the `generate` command. Customize the length, character sets, and exclusions as needed.

### Exporting and Importing Data

```sh
# Export data
passman export <file> -f <format>

# Import data
passman import <file> -f <format>
```

Export and import data from and to the vault in various formats like CSV and JSON.

### Changing the Master Password

```sh
passman change-master
```

Change the master password of the vault with the `change-master` command.

### Configuring Passman

```sh
passman config
```

Configure Passman settings using the `config` command.

### Locking and Unlocking the Vault

```sh
# Lock the vault
passman lock

# Unlock the vault
passman unlock
```

Manually lock and unlock the password vault for added security.

## Options

- `-h, --help`: Show the help message and exit.

- `--version`: Show the program's version number and exit.

For each command, specific options are available as mentioned in the usage examples above.