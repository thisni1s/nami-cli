# NaMi CLI

nami-cli allows you to interact with the DPSG NaMi and get information about your Members.

It is written in Go and uses the [nami-go](https://github.com/thisni1s/nami-go) library for interacting with the NaMi and [Cobra](https://github.com/spf13/cobra) for building the CLI Application. 

## Configuration

A config file containing your credentials is needed for this to work.

By default the file should be located at ```.nami.yml``` and look like this:
```
username: 133337 # your nami id
password: verysecure
gruppierung: 010101 # your "Stammesnummer"
```

## Usage
```
nami-cli [command]

nami-cli info 133337                 # Prints info about Member with ID 133337
nami-cli search name -f John -l Doe  # Searches for Members with first name "John" and last name "Doe"
nami-cli search occupation leiter    # Searches for Members with occupation (Tätigkeit) "Leiter"
nami-cli search subdivision rover    # Searches for Members with subdivision (Untergliederung) "Rover"
nami-cli search tag 1337             # Searches for Members with TagId "1337"

```

### Available Commands:
- info        Prints information about a specified Member
- search      Search for different kinds of Members in Nami
- completion  Generate the autocompletion script for the specified shell
- help        Help about any command

### Info
Prints information about a user specified by their Member ID.
The output is YAML but can be switched to indented JSON with the --json flag.

#### Flags:
- -h, --help   help for info
- --json   Print the Info as JSON.

### Search
Search Nami for Members visible to the logged in User.
Different filters are provided with the use of subcommands.

#### Available Sub-Commands:
- name        Search for Members by Name
- occupation  Search for Members with a specific occupation
- subdivision Search for members in a specific subdivision
- tag         Search Members by Tag

#### Flags:
- -e, --email   Output found members in mailbox format e.g. 'John Doe <john@example.com>' (only prints members that have a mail address!!)  
- -h, --help    help for search

### Global Flags:
- --config string   config file (default is .nami.yaml)
- -h, --help            help for nami-cli




