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

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag [tagID]",
	Short: "Search Members by Tag",
	Long: `Search for Members using the ID of a Tag.

Example:
  nami-cli search tag 1337`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		tag := args[0]
		//tag, err := strconv.Atoi(args[0])
		//if err != nil {
		//    log.Printf("Tag must be number!")
		//    log.Fatal(err)
		//}
        Login()
		list, err := namigo.Search(namiTypes.SearchValues{
			TagID: tag,
		})
		if err != nil {
			log.Println("Failed to get Members for provided Tag!")
			log.Fatal(err)
		}
        PrintSearchResult(list)

	},
}

func init() {
	searchCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
