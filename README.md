# NaMi CLI

nami-cli allows you to interact with the DPSG NaMi and get information about your Members.

It is written in Go and uses the [nami-go](https://github.com/thisni1s/nami-go) library for interacting with the NaMi and [Cobra](https://github.com/spf13/cobra) for building the CLI Application. 

## Configuration

A config file containing your credentials is needed for this to work.

By default the file should be located at ```~/.nami.yml``` and look like this:
```
username: 133337 # your nami id
password: verysecure
gruppierung: 010101 # your "Stammesnummer"
```
Different Commands may neeed additional config files!


## Usage
```
nami-cli [command]

nami-cli info 133337                 # Prints info about Member with ID 133337
nami-cli search name John Doe        # Searches for Members with first name "John" and last name "Doe"
nami-cli search occupation leiter    # Searches for Members with occupation (Tätigkeit) "Leiter"
nami-cli search subdivision rover    # Searches for Members with subdivision (Untergliederung) "Rover"
nami-cli search tag 1337             # Searches for Members with TagId "1337"
nami-cli search -n John -l Doe -o leiter -d rover -t 1337    # Combine all above filters into one using flags
nami-cli mail -t 1337 --mailCfg mailconfig.yml # Send E-Mail to all members with tag 1337
nami-cli sepa --all --fee 20.00 --out res.xml # Generate SEPA XML File for all members for 20.00€

```

### Available Commands:
- info        Prints information about a specified Member
- search      Search for different kinds of Members in Nami
- mail        Send E-Mails to different Members
- sepa        Generate SEPA XML files.
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
Different filters are provided with the use of subcommands or flags.  
For possible flag values consult the help command of the specific sub command.  
Normal Output is of the form ```ID: FirstName LastName``` but can be changed to mailbox, YAML or JSON format
with the use of the ```--email``` ```--full``` or ```--json``` flags.

#### Available Sub-Commands:
- name        Search for Members by Name
- occupation  Search for Members with a specific occupation
- subdivision Search for members in a specific subdivision
- tag         Search Members by Tag

#### Flags:
- -n, --fname string         First name (if any)
- -l, --lname string         Last name (if any)
- -o, --occupation string    Occupation (if any) for options see 'occupation' sub command help
- -d, --subdivision string   Subdivision (if any) for options see 'subdivision' sub command help
- -t, --tag string           Tag (if any)
- -e, --email                Output found members in mailbox format e.g. 'John Doe <john@example.com>' (only prints members that have a mail address!!) 
- -f, --full                 Fully output found members (in YAML format)
- -j, --json                 Output found members in JSON format
- -h, --help                 help for search

### Sepa
Generate SEPA XML file for specified users
You can specify users with the ```--tag```, ```--occupation```, ```--subdivision``` and ```--all``` tags just like the search.
Fixed fees can be set with ```--fee```.
Output file needs to be specified with ```--out```
A special SEPA config file is needed! Location can be specified with ```--sepaConfig``` See ```sepa.yml.example``` for an example!

#### Flags
- -a, --all                  Create file for ALL members
- --fee float            Fixed Fee, ignore member fees and set a fixed fee
- -h, --help                 help for sepa
- -o, --occupation string    Occupation (if any) for options see 'occupation' sub command help
- --out string           Output file
- -s, --sepaConfig string    Path to the sepa config, default is ~/sepa.yml
- -d, --subdivision string   Subdivision (if any) for options see 'subdivision' sub command help
- -t, --tag string           Tag (if any)

### Mail
Send E-Mails to different Members.
Specify whom to send the E-Mails to using the flags!
E-Mail content can be defined with a template file. 
In it you have access to all Fields of a Member, plus their Beitrag with ```.FixBeitrag```
Specify everything related to the E-Mail in the ```mailconfig.yml``` file!
A special Mail config file is needed! Location can be specified with ```--mailCfg``` See ```.mailconfig.yml.example``` for an example!
A special Template file is needed! Location can be specified within the ```mailconfig.yml``` See ```message.tmpl.example``` for an example!

#### Flags
- -a, --all                  Send E-Mail to ALL members
- -h, --help                 help for mail
- --mailCfg string       E-Mail config file. Defaults to ~/.mailconfig.yml
- -o, --occupation string    Occupation (if any) for options see 'occupation' sub command help
- -d, --subdivision string   Subdivision (if any) for options see 'subdivision' sub command help
- -t, --tag string           Tag (if any)



### Global Flags:
- --config string   config file (default is ~/.nami.yaml)
- -h, --help            help for nami-cli





