package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"gopkg.in/gomail.v2"
//	"gopkg.in/yaml.v3"
	"log"
//	"os"
	"text/template"
)

// sepaCmd represents the sepa command
var sepaCmd = &cobra.Command{
	Use:   "sepa",
	Short: "Send a SEPA Pre Notification",
	Long: `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		Login()
		list, err := namigo.Search(namiTypes.SearchValues{
			TagID: *MailTag,
		})
		log.Println(list)

		if err != nil || len(list) == 0 {
			log.Println("Failed to get Members for provided Tag!")
			log.Fatal(err)
		}

		var sepalist []SepaMailbox
		for _, member := range list {
			log.Println("tryna get info for: ", member.Vorname)
			fullMember, err := namigo.GetMemberDetails(fmt.Sprint(member.EntriesID), config.Gruppierung)
			if err != nil {
				log.Println("Failed to get Member information!")
				log.Fatal(err)
			}
			var mail string
			if fullMember.Email != "" {
				mail = fullMember.Email
			} else if fullMember.EmailVertretungsberechtigter != "" {
				mail = fullMember.EmailVertretungsberechtigter
			} else {
				log.Fatal("No Email Address found!")
			}
			sepamem := SepaMailbox{
				Name:      fullMember.Kontoverbindung.Kontoinhaber,
				ChildName: fullMember.Vorname + " " + fullMember.Nachname,
				Address:   mail,
				Amount:    calcPaymentAmount(fullMember),
				Reference: fmt.Sprintf("%d-%s-%s", fullMember.MitgliedsNummer, fullMember.Vorname, fullMember.Nachname),
			}
			sepalist = append(sepalist, sepamem)
		}
		sendMail(sepalist)

	},
}

func calcPaymentAmount(member namiTypes.Member) string {
	if member.Beitragsart == "Voller Beitrag" || member.Beitragsart == "Voller Beitrag - Stiftungseuro" {
		return mailCfg.FeeFull
	} else if member.Beitragsart == "Familienermäßigt" || member.Beitragsart == "Familienermäßigt - Stiftungseuro" {
		return mailCfg.FeeFam
	} else if member.Beitragsart == "Sozialermäßigt" || member.Beitragsart == "Sozialermäßigt - Stiftungseuro" {
		return mailCfg.FeeSocial
	} else {
		return mailCfg.FeeFull
	}
}

type SepaMailConfig struct {
	Username  string
	Fromname  string
	Password  string
	Server    string
	Subject   string
	Template  string
	FeeFull   string
	FeeFam    string
	FeeSocial string
}

var mailCfg SepaMailConfig

func init() {
	mailCmd.AddCommand(sepaCmd)

	var cfgHandle string
	mailCmd.Flags().StringVar(&cfgHandle, "mailCfg", "", "E-Mail config file")

    /*
	file, err := os.ReadFile(cfgHandle)
	if err != nil {
		log.Println("Failed to read E-Mail config file!!")
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &mailCfg)
	if err != nil {
		log.Println("Failed to read E-Mail config file!!")
		log.Fatal(err)
	}
    */

}

func sendMail(list []SepaMailbox) {
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

	for _, r := range list {
		m.SetHeader("From", mailCfg.Username)
		m.SetHeader("To", r.Address)
		m.SetAddressHeader("Cc", mailCfg.Username, mailCfg.Fromname)
		m.SetHeader("Subject", mailCfg.Subject)

		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, r)
		if err != nil {
			log.Fatal("Failed to Execute template")
		}

		m.SetBody("text/plain", buf.String())

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r.Address, err)
		}
		m.Reset()
	}

}
