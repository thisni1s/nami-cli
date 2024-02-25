package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	"gopkg.in/yaml.v3"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Prints information about a specified Member",
	Long: `Prints information about a user specified by their Member ID.
The output is YAML but can be switched to indented JSON with the --json flag.

Examples:
  nami-cli info 133337
  nami-cli info 133337 --json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        readConfig()
		Login()
		mem, err := namigo.GetMemberDetails(args[0], GetGroupId())
		if err != nil {
			log.Println("Error retrieving info about user!")
			log.Fatal(err)
		}
		var s []byte
		if !*asJson {
			s, _ = yaml.Marshal(mem)
		} else {
			s, _ = json.MarshalIndent(mem, "", "    ")
		}
		fmt.Printf("%s \n", s)
	},
}

var asJson *bool

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	asJson = infoCmd.Flags().Bool("json", false, "Print the Info as JSON.")
}
