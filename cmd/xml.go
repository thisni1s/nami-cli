/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"time"

	"github.com/almerlucke/go-iban/iban"
	"github.com/apsl/sepakit/sepadebit"
	"github.com/spf13/cobra"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"gopkg.in/yaml.v3"
)

// xmlCmd represents the xml command
var xmlCmd = &cobra.Command{
	Use:   "xml",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		list := findMembers()
		fullList := getMemberDetails(list)
		genSepaXml(fullList)

	},
}

func findMembers() *[]namiTypes.SearchMember {
	log.Println("Finding members")
	var sValues namiTypes.SearchValues
	if *occupation != "" {
		sValues.TaetigkeitID = CheckOccupationArg(*occupation)
	}
	if *subdivision != "" {
		sValues.UntergliederungID = CheckSubdivisionArg(*subdivision)
	}
	if *tag != "" {
		sValues.TagID = *tag
	}
	Login()
	list, err := namigo.Search(sValues)
	if err != nil {
		log.Println("Something went wrong searching for Members!")
		log.Fatal(err)
	}
	return &list
}

func getMemberDetails(searchres *[]namiTypes.SearchMember) *[]namiTypes.Member {
	log.Println("Getting member details")
	var wg sync.WaitGroup
	var res = make([]namiTypes.Member, len(*searchres))
	for i, mem := range *searchres {
		wg.Add(1)
		go func(i int, mem namiTypes.SearchMember) {
			defer wg.Done()
			id := strconv.Itoa(mem.ID)
			member, err := namigo.GetMemberDetails(id, config.Gruppierung)
			if err != nil {
				log.Println("Error getting details of: " + mem.Vorname + mem.Nachname)
				log.Println(err)
			}
			//1999-03-24 00:00:00
			age, _ := time.Parse("2006-01-02 03:04:05", member.GeburtsDatum)
			now := time.Now()
			if now.Sub(age).Hours() >= 18*365*24 { // is member over 18?
				activs, err := namigo.GetMemberActivities(id)
				if err != nil {
					log.Printf("error getting activities for member: %s %s \n", mem.Vorname, mem.Nachname)
				}
				var leader bool
				for _, act := range activs {
					if act.Taetigkeit == "€ LeiterIn (6)" {
						leader = true
					}
				}
				if leader {
					member.GenericField1 = "leader"
				}
			}
			res[i] = member
		}(i, mem)
	}
	wg.Wait()
	return &res
}

func genSepaXml(list *[]namiTypes.Member) {

	var ctrlSum float64

	docxml := sepadebit.NewDocument()
	docxml.SetInitiatingParty(sepaCfg.CreditorName, sepaCfg.CreditorID)
	docxml.SetCreationDateTime(time.Now())

	paym := sepadebit.NewPayment()

	paym.Creditor = sepadebit.NewCreditor()
	paym.Creditor.ID = sepaCfg.CreditorID
	paym.Creditor.BIC = sepaCfg.CreditorBIC
	paym.Creditor.IBAN = sepaCfg.CreditorIBAN
	paym.Creditor.Name = sepaCfg.CreditorName
	paym.Creditor.ChargeBearer = "SLEV"

	paym.ID = fmt.Sprintf("%s-%s", strings.ReplaceAll(sepaCfg.CreditorName, " ", ""), time.Now().Format("20060102"))
	paym.Method = "DD"
	paym.RequestedCollectionDate = sepaCfg.CollectionDate

	for _, member := range *list {
		valid := verifyMandate(member)
        if !valid {
            log.Printf("Mandate verification failed for: %s %s \n", member.Vorname, member.Nachname)
        }
		date, err := time.Parse("02012006", member.Kontoverbindung.Kontonummer)
		if err != nil {
			valid = false
			log.Printf("Error with Mandate Date of: %s %s \n", member.Vorname, member.Nachname)
			log.Println(err)
		}
        amount := calcAmount(member.BeitragsartID, member.GenericField1, member.MglType)
        if amount == 0.00 {
            log.Printf("WRONG Amount for: %s %s \n", member.Vorname, member.Nachname)
            valid = false
        }
		if valid {
			ctrlSum += amount

			tr := sepadebit.Transaction{
				ID:             "NOTPROVIDED",
				MandateID:      fmt.Sprintf("%d-%s-%s", member.MitgliedsNummer, strings.ReplaceAll(member.Vorname, " ", ""), member.Nachname),
				Date:           sepadebit.Date(date),
				RemittanceInfo: sepaCfg.CollectionInfo,
				Debtor: sepadebit.Debtor{
					Name: member.Kontoverbindung.Kontoinhaber,
					BIC:  member.Kontoverbindung.Bic,
					IBAN: member.Kontoverbindung.Iban,
				},
				Amount: sepadebit.TAmount{
					Amount:   f2s(amount),
					Currency: "EUR",
				},
			}
			paym.Transactions = append(paym.Transactions, tr)

		}

	}

	paym.TransacNb = len(paym.Transactions)
	paym.CtrlSum = f2s(ctrlSum)

	docxml.AddPayment(paym)
	docxml.TransacNb = len(paym.Transactions)
	docxml.CtrlSum = f2s(ctrlSum)

	file, err := os.Create(*outFile)
	if err != nil {
		log.Fatal("cant create file")
	}
	err = docxml.WriteLatin1(file)
	if err != nil {
		log.Fatal(err)
	}

}

func verifyMandate(member namiTypes.Member) bool {
	if member.Kontoverbindung.Kontoinhaber == "" {
		return false
	} else if member.Kontoverbindung.Kontonummer == "" {
		return false
	} else if member.Kontoverbindung.Bic == "" {
		return false
	}
	_, err := iban.NewIBAN(member.Kontoverbindung.Iban)
	if err != nil {
		return false
	}
	return true
}

func calcAmount(beitragsart int, taetigkeit string, mgltype string) float64 {
	if *fixedFee != 0.00 {
		return *fixedFee
	} else if taetigkeit == "leader" {
		return sepaCfg.FeeLeader
	} else if beitragsart == 1 || beitragsart == 4 {
		return sepaCfg.FeeFull
	} else if beitragsart == 2 || beitragsart == 5 {
		return sepaCfg.FeeFamily
	} else if beitragsart == 3 || beitragsart == 6 {
		return sepaCfg.FeeSocial
	} else if mgltype == "MITGLIED" {
		return sepaCfg.FeePassive
	} else {
		return 0.00
	}
}

func f2s(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

type SepaConfig struct {
	CreditorID   string
	CreditorName string
	CreditorIBAN string
	CreditorBIC  string

	CollectionDate string
	CollectionInfo string

	FeeFull    float64
	FeeFamily  float64
	FeeSocial  float64
	FeeLeader  float64
	FeePassive float64
}

var fixedFee *float64
var sepaCfg SepaConfig
var outFile *string
var searchAll *bool

func init() {
	rootCmd.AddCommand(xmlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xmlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// xmlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	occupation = xmlCmd.Flags().StringP("occupation", "o", "", "Occupation (if any) for options see 'occupation' sub command help")
	subdivision = xmlCmd.Flags().StringP("subdivision", "d", "", "Subdivision (if any) for options see 'subdivision' sub command help")
	tag = xmlCmd.Flags().StringP("tag", "t", "", "Tag (if any)")
	searchAll = xmlCmd.Flags().BoolP("all", "a", false, "Create file for ALL members")
	xmlCmd.MarkFlagsOneRequired("occupation", "subdivision", "tag", "all")
	xmlCmd.MarkFlagsMutuallyExclusive("all", "occupation")
	xmlCmd.MarkFlagsMutuallyExclusive("all", "subdivision")
	xmlCmd.MarkFlagsMutuallyExclusive("all", "tag")

	fixedFee = xmlCmd.Flags().Float64("fee", 0.00, "Fixed Fee, ignore member fees and set a fixed fee")

	var cfgFile *string
	cfgFile = xmlCmd.Flags().StringP("sepaConfig", "s", "", "Path to the sepa config, default is ~/sepa.yml")
	if *cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Error finding your home directory!")
		}
		*cfgFile = fmt.Sprintf("%s/sepa.yml", home)
	}

	file, err := os.ReadFile(*cfgFile)
	if err != nil {
		log.Println("Failed to read sepa config file!")
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &sepaCfg)
	if err != nil {
		log.Println("Failed to read sepa config file!")
		log.Fatal(err)
	}

	outFile = xmlCmd.Flags().String("out", "", "Output file")
	xmlCmd.MarkFlagRequired("out")

}
