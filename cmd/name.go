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
This command supports two arguments, if one argument is given nami-cli will search for first and last names mathing this argument.
If two arguments are given, it will search for a Member matching both arguments

Examples:
  nami-cli search name John
  nami-cli search name Doe
  nami-cli search name John Doe`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var searchCfg namiTypes.SearchValues
		Login()
		if len(args) == 1 {
			searchCfg.Vorname = args[0]
			list, err := namigo.Search(searchCfg)
			if err != nil {
				log.Println("Failed to get Members for provided first name!")
				log.Fatal(err)
			}
			searchCfg.Vorname = ""
			searchCfg.Nachname = args[0]
			list2, err := namigo.Search(searchCfg)
			if err != nil {
				log.Println("Failed to get Members for provided last name!")
				log.Fatal(err)
			}
			PrintSearchResult(append(list, list2...))
		} else if len(args) == 2 {
            searchCfg.Vorname = args[0]
            searchCfg.Nachname = args[1]
			list, err := namigo.Search(searchCfg)
			if err != nil {
				log.Println("Failed to get Members for provided Name!")
				log.Fatal(err)
			}
			PrintSearchResult(list)
		}
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
	//fName = nameCmd.Flags().StringP("fName", "v", "", "First name (if any)")
	//lName = nameCmd.Flags().StringP("lName", "l", "", "Last name (if any)")
	//nameCmd.MarkFlagsOneRequired("fName", "lName")
}
