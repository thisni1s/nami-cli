package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	namiTypes "github.com/thisni1s/nami-go/types"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"text/template"
)

// mailCmd represents the mail command
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Send E-Mails to different Members",
	Long: `Send E-Mails to different Members.
Specify whom to send the E-Mails to using the flags!
E-Mail content can be defined with a template file. 
In it you have access to all Fields of a Member, plus their Beitrag with .FixBeitrag
Specify everything related to the E-Mail in the mailconfig.yml file!

Examples:
  nami-cli mail --tag 1337
  nami-cli mail --mailCfg /home/user/mailconfig.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		cpFlags()
		readMailCfg()
		list := findMembers()
		fullList := getMemberDetails(list)
		cleanList := cleanAddresses(fullList)
		sendMail(cleanList)

	},
}

var mailCfgHandle *string
var mailCfg SepaMailConfig
var mailTag *string
var mailOcc *string
var mailSub *string
var mailSa *bool

func sendMail(list *[]namiTypes.Member) {
	d := gomail.NewDialer(mailCfg.Server, 587, mailCfg.Username, mailCfg.Password)
	s, err := d.Dial()
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	m := gomail.NewMessage()
	tmpl, err := template.New(mailCfg.Template).ParseFiles(mailCfg.Template)
	if err != nil {
		log.Fatal("Failed to read template ", err)
	}

	for _, r := range *list {
		m.SetHeader("From", mailCfg.Username)
		m.SetHeader("To", r.Email)
		m.SetAddressHeader("Cc", mailCfg.Username, mailCfg.Fromname)
		m.SetHeader("Subject", mailCfg.Subject)

		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, r)
		if err != nil {
			log.Println("Failed to Execute template")
			log.Fatal(err)
		}

		m.SetBody("text/plain", buf.String())

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r.Email, err)
		}
		m.Reset()
	}

}

func cleanAddresses(list *[]namiTypes.Member) *[]namiTypes.Member {
	var res []namiTypes.Member
	for _, mem := range *list {
		if mem.Email == "" {
			if mem.EmailVertretungsberechtigter != "" {
				mem.Email = mem.EmailVertretungsberechtigter
				mem.FixBeitrag = calcMailAmount(mem.BeitragsartID, mem.GenericField1, mem.MglType)
				res = append(res, mem)
			} else {
				fmt.Printf("No E-Mail Address found for: %s %s Not sending any mail to them! \n", mem.Vorname, mem.Nachname)
			}
		} else {
			mem.FixBeitrag = calcMailAmount(mem.BeitragsartID, mem.GenericField1, mem.MglType)
			res = append(res, mem)
		}
	}
	return &res
}

func calcMailAmount(beitragsart int, taetigkeit string, mgltype string) float64 {
	if taetigkeit == "leader" {
		return mailCfg.FeeLeader
	} else if beitragsart == 1 || beitragsart == 4 {
		return mailCfg.FeeFull
	} else if beitragsart == 2 || beitragsart == 5 {
		return mailCfg.FeeFam
	} else if beitragsart == 3 || beitragsart == 6 {
		return mailCfg.FeeSocial
	} else if mgltype == "MITGLIED" {
		return mailCfg.FeePassive
	} else {
		return 0.00
	}
}

// no idea why this is neccessary but it does not work without it
func cpFlags() {
	*tag = *mailTag
	*occupation = *mailOcc
	*subdivision = *mailSub
	*searchAll = *mailSa
}

type SepaMailConfig struct {
	Username   string  `yaml:"Username"`
	Fromname   string  `yaml:"Fromname"`
	Password   string  `yaml:"Password"`
	Server     string  `yaml:"Server"`
	Subject    string  `yaml:"Subject"`
	Template   string  `yaml:"Template"`
	FeeFull    float64 `yaml:"FeeFull"`
	FeeFam     float64 `yaml:"FeeFam"`
	FeeSocial  float64 `yaml:"FeeSocial"`
	FeeLeader  float64 `yaml:"FeeLeader"`
	FeePassive float64 `yaml:"FeePassive"`
}

func readMailCfg() {
	if *mailCfgHandle == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Error finding your home directory!")
		}
		*mailCfgHandle = fmt.Sprintf("%s/.mailconfig.yml", home)
	}
	file, err := os.ReadFile(*mailCfgHandle)
	if err != nil {
		log.Println("Failed to read mail config file!")
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &mailCfg)
	if err != nil {
		log.Println("Failed to read mail config file!")
		log.Fatal(err)
	}

}

func init() {
	rootCmd.AddCommand(mailCmd)

	mailOcc = mailCmd.Flags().StringP("occupation", "o", "", "Occupation (if any) for options see 'occupation' sub command help")
	mailSub = mailCmd.Flags().StringP("subdivision", "d", "", "Subdivision (if any) for options see 'subdivision' sub command help")
	mailTag = mailCmd.Flags().StringP("tag", "t", "", "Tag (if any)")
	mailSa = mailCmd.Flags().BoolP("all", "a", false, "Send E-Mail to ALL members")

	mailCmd.MarkFlagsOneRequired("occupation", "subdivision", "tag", "all")
	mailCmd.MarkFlagsMutuallyExclusive("all", "occupation")
	mailCmd.MarkFlagsMutuallyExclusive("all", "subdivision")
	mailCmd.MarkFlagsMutuallyExclusive("all", "tag")

	mailCfgHandle = mailCmd.Flags().String("mailCfg", "", "E-Mail config file. Defaults to ~/.mailconfig.yml")

}
