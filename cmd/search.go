package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"gopkg.in/yaml.v3"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for different kinds of Members in Nami",
	Long: `Search Nami for Members visible to the logged in User.
    Different filters are provided with the use of subcommands.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		readConfig()
		var sValues namiTypes.SearchValues
		if firstname != "" {
			sValues.Vorname = firstname
		}
		if lastname != "" {
			sValues.Nachname = lastname
		}
		if occupation != "" {
			sValues.TaetigkeitID = CheckOccupationArg(occupation)
		}
		if subdivision != "" {
			sValues.UntergliederungID = CheckSubdivisionArg(subdivision)
		}
		if tag != "" {
			sValues.TagID = tag
		}
        addMemberTypes(&sValues)
		Login()
		list, err := namigo.Search(sValues)
		if err != nil {
			log.Println("Something went wrong searching for Members!")
			log.Fatal(err)
		}
		PrintSearchResult(list)
	},
}

var email bool
var jsono bool
var fullo bool
var nonMembers bool
var schnupper bool
var inactive bool

var firstname string
var lastname string
var occupation string
var subdivision string
var tag string

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	searchCmd.PersistentFlags().BoolVarP(&email, "email", "e", false, "Output found members in mailbox format e.g. 'John Doe <john@example.com>' (only prints members that have a mail address!!) ")
	searchCmd.PersistentFlags().BoolVarP(&jsono, "json", "j", false, "Output found members in JSON format")
	searchCmd.PersistentFlags().BoolVarP(&fullo, "full", "f", false, "Fully output found members (in YAML format)")
	searchCmd.PersistentFlags().BoolVar(&nonMembers, "nonMembers", false, "Also search for non Members (NICHT_MITGLIED)")
	searchCmd.PersistentFlags().BoolVar(&schnupper, "schnupper", false, "Also search for Schnupper Members (SCHNUPPER_MITGLIED)")
	searchCmd.PersistentFlags().BoolVar(&inactive ,"inactive", false, "Search for inactive Members (INAKTIV)")
	searchCmd.MarkFlagsMutuallyExclusive("email", "json", "full")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//firstname= searchCmd.Flags().StringP("fname", "n", "", "First name (if any)")
	searchCmd.Flags().StringVarP(&firstname, "fname", "n", "", "First name (if any)")
	searchCmd.Flags().StringVarP(&lastname, "lname", "l", "", "Last name (if any)")
	searchCmd.Flags().StringVarP(&occupation, "occupation", "o", "", "Occupation (if any) for options see 'occupation' sub command help")
	searchCmd.Flags().StringVarP(&subdivision, "subdivision", "d", "", "Subdivision (if any) for options see 'subdivision' sub command help")
	searchCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag (if any)")
	searchCmd.MarkFlagsOneRequired("fname", "lname", "occupation", "subdivision", "tag")
}

func addMemberTypes(sValues *namiTypes.SearchValues) {
    sValues.MglTypeID = namiTypes.MITGLIED
    if nonMembers == true {
        sValues.MglTypeID = namiTypes.NICHTMITGLIED
    }
    if schnupper == true {
        sValues.MglTypeID = namiTypes.SCHNUPPER
    }
    if inactive == false {
        sValues.MglStatusID = namiTypes.AKTIV
    } else {
        sValues.MglStatusID = namiTypes.INAKTIV
    }
}

func PrintSearchResult(members []namiTypes.SearchMember) {
	if email {
		for _, mem := range members {
			if mem.Email != "" || mem.EmailVertretungsberechtigter != "" {
				var mail string
				if mem.Email != "" {
					mail = mem.Email
				} else {
					mail = mem.EmailVertretungsberechtigter
				}
				fmt.Printf("%s %s <%s> \n", mem.Vorname, mem.Nachname, mail)
			}
		}
	} else if jsono {
		s, _ := json.MarshalIndent(members, "", "    ")
		fmt.Printf("%s, \n", s)
	} else if fullo {
		s, _ := yaml.Marshal(members)
		fmt.Printf("%s, \n", s)
	} else {
		for _, mem := range members {
			fmt.Printf("%d: %s %s \n", mem.ID, mem.Vorname, mem.Nachname)
		}
	}
}

func CheckOccupationArg(input string) namiTypes.TAETIGKEIT {
	var ugId namiTypes.TAETIGKEIT
	switch input {
	case "mgl":
		ugId = namiTypes.TG_MITGLIED
	case "elternv":
		ugId = namiTypes.TG_ELTERNVERTRETUNG
	case "leiter":
		ugId = namiTypes.TG_LEITER
	case "delegiert":
		ugId = namiTypes.TG_DELEGIERT
	case "beobachter":
		ugId = namiTypes.TG_BEOBACHTER
	case "kurat":
		ugId = namiTypes.TG_KURAT
	case "vorsitz":
		ugId = namiTypes.TG_VORSITZ
	case "admin":
		ugId = namiTypes.TG_ADMIN
	case "sonstmit":
		ugId = namiTypes.TG_SONSTMITARBEITER
	case "akmgl":
		ugId = namiTypes.TG_MITGLIEDAK
	case "gf":
		ugId = namiTypes.TG_GESCHÄFTSFÜHRER
	case "kassierer":
		ugId = namiTypes.TG_KASSIERER
	case "prüfer":
		ugId = namiTypes.TG_KASSENPRÜFER
	case "matwart":
		ugId = namiTypes.TG_MATWART
	case "schnupper":
		ugId = namiTypes.TG_SCHNUPPER
	case "passiv":
		ugId = namiTypes.TG_PASSIV
	case "sonstmgl":
		ugId = namiTypes.TG_SONSTMITGLIED
	case "sonstext":
		ugId = namiTypes.TG_SONSTEXT
	default:
		log.Fatal("You need to provide an occupation!")
	}
	return ugId
}

func CheckSubdivisionArg(input string) namiTypes.UNTERGLIEDERUNG {
	var ugId namiTypes.UNTERGLIEDERUNG
	switch input {
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
    case "biber":
        ugId = namiTypes.UG_BIBER
	default:
		log.Fatal("You need to provide a subdivison!")
	}
	return ugId
}
