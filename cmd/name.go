/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
)

// nameCmd represents the name command
var nameCmd = &cobra.Command{
	Use:   "name",
	Short: "Search for Members by Name",
	Long: `Returns all Members whos name match the specified name.
Both fName and lName are optional, but one is required

Example:
  nami-cli search name --fName John --lName Doe`,
	Run: func(cmd *cobra.Command, args []string) {
        if *fName == "" && *lName == "" {
            log.Fatal("At least one Flag is required!")
        }
		Login()
		var searchCfg namiTypes.SearchValues
		if *fName != "" {
			searchCfg.Vorname = *fName
		}
		if *lName != "" {
			searchCfg.Nachname = *lName
		}
		list, err := namigo.Search(searchCfg)
		if err != nil {
			log.Println("Failed to get Members for provided Name!")
			log.Fatal(err)
		}
		PrintSearchResult(list)

	},
}

var fName *string
var lName *string

func init() {
	searchCmd.AddCommand(nameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	fName = nameCmd.Flags().StringP("fName", "v", "", "First name (if any)")
	lName = nameCmd.Flags().StringP("lName", "l", "", "Last name (if any)")
    nameCmd.MarkFlagsOneRequired("fName", "lName")
}
