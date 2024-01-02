/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	namiTypes "github.com/thisni1s/nami-go/types"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for different kinds of Members in Nami",
	Long: `Search Nami for Members visible to the logged in User.
    Different filters are provided with the use of subcommands.
    `,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println(Config.Username)
	//},
}

var email *bool

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	email = searchCmd.PersistentFlags().BoolP("email", "e", false, "Output found members in mailbox format e.g. 'John Doe <john@example.com>' (only prints members that have a mail address!!) ")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func PrintSearchResult(members []namiTypes.SearchMember) {
	for _, mem := range members {
		if *email {
			if mem.Email != "" {
				fmt.Printf("%s %s <%s> \n", mem.Vorname, mem.Nachname, mem.Email)
			} 
		} else {
            fmt.Printf("%d: %s %s \n", mem.ID, mem.Vorname, mem.Nachname)
        }
	}

}
