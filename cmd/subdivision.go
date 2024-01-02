/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"log"
)

// subdivisionCmd represents the subdivision command
var subdivisionCmd = &cobra.Command{
	Use:       "subdivision [woe|juffi|pfadi|rover|stavo|sonst]",
	Short:     "Search for members in a specific subdivision",
	Long:      `Prints all members in the specified subdivision
    Possible subdivions are: "woe", "juffi", "pfadi", "rover", "stavo", "sonst"

Example:
  nami-cli search subdivision rover`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"woe", "juffi", "pfadi", "rover", "stavo", "sonst"},
	Run: func(cmd *cobra.Command, args []string) {
		var ugId namiTypes.UNTERGLIEDERUNG
		switch args[0] {
		case "woe":
			ugId = namiTypes.UG_WOE
		case "juffi":
			ugId = namiTypes.UG_JUFFI
		case "pfadi":
			ugId = namiTypes.UG_PFADI
		case "rover":
			ugId = namiTypes.UG_ROVER
		case "stavo":
			ugId = namiTypes.UG_STAVO
		case "sonst":
			ugId = namiTypes.UG_SONST
		default:
            log.Fatal("You need to provide a subdivison!")
		}

		Login()

		list, err := namigo.Search(namiTypes.SearchValues{
			UntergliederungID: ugId,
		})
		if err != nil {
			log.Println("Failed to get Members for provided Tag!")
			log.Fatal(err)
		}
		PrintSearchResult(list)

	},
}

func init() {
	searchCmd.AddCommand(subdivisionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subdivisionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subdivisionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
