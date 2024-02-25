package cmd

import (
	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"log"
)

// occupationCmd represents the occupation command
var occupationCmd = &cobra.Command{
	Use:   "occupation [occupationShorthand]",
	Short: "Search for Members with a specific occupation",
	Long: `Prints all Members that have a specified occupation.
Supported occupations are:
  Mitglied -> mgl
  Schnupper Mitglieder -> schnupper
  Passive Mitglieder -> passiv
  Leiter -> leiter
  Kuraten -> kurat
  Vorsitz -> vorsitz
  Kassierer -> kassierer
  Kassenprüfer -> prüfer
  Materialwart -> matwart
  Elternvertreter -> elternv
  Admin -> admin
  AK Mitglieder -> akmgl
  Geschaftsführer -> gf
  Beobachter -> beobachter
  Delegierte -> delegiert
  Sonstige Mitarbeiter -> sonstmit
  Sonstige Mitglieder -> sonstmgl
  Sonstige Externe -> sonstext

Examples:
  nami-cli search occupation leiter
  nami-cli search occupation mgl`,
	ValidArgs: []string{"mgl", "elternv", "leiter", "delegiert", "beobachter", "kurat", "vorsitz", "admin", "sonstmit", "akmgl", "gf", "kassierer", "prüfer", "matwart", "schnupper", "passiv", "sonstmgl", "sonstext"},
	Run: func(cmd *cobra.Command, args []string) {
        readConfig()
		ugId := CheckOccupationArg(args[0])
		Login()
        sValues := namiTypes.SearchValues{
            TaetigkeitID: ugId,
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
	searchCmd.AddCommand(occupationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// occupationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// occupationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
