/*
Copyright © 2024 thisni1s thisni1s@jn2p.de

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"log"
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
		readConfig()
		tag := args[0]
		//tag, err := strconv.Atoi(args[0])
		//if err != nil {
		//    log.Printf("Tag must be number!")
		//    log.Fatal(err)
		//}
		Login()
		sValues := namiTypes.SearchValues{
			TagID: tag,
		}
        addMemberTypes(&sValues)
		list, err := namigo.Search(sValues)
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
