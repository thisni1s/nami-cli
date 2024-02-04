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
	"github.com/thisni1s/nami-cli/helpers"
	namigo "github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"gopkg.in/yaml.v3"
)

// sepaCmd represents the sepa command
var sepaCmd = &cobra.Command{
	Use:   "sepa",
	Short: "Generate SEPA XML files.",
	Long: `Generate SEPA XML file for specified users
You can specify users with the --tag, --occupation, --subdivision and --all tags just like the search.
Fixed fees can be set with --fee.
Output file needs to be specified with --out
A special SEPA config file is needed! Location can be specified with --sepaConfig

Examples:
  nami-cli sepa -d rover --out result.xml
  nami-cli sepa --all --fee 20.00 --out result.xml
`,
	Run: func(cmd *cobra.Command, args []string) {
		readSepaCfg()
		list := findMembers()
		fullList := getMemberDetails(list)
		genSepaXml(fullList)

	},
}

func findMembers() *[]namiTypes.SearchMember {
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
	fmt.Println("Finding members")
	list, err := namigo.Search(sValues)
	if err != nil {
		log.Println("Something went wrong searching for Members!")
		log.Fatal(err)
	}
	return &list
}

func getMemberDetails(searchres *[]namiTypes.SearchMember) *[]namiTypes.Member {
	var wg sync.WaitGroup
	var res = make([]namiTypes.Member, len(*searchres))

	progress := make(chan int)
	wg.Add(1)
	go printProgress(&progress, len(*searchres), &wg)

	for i, mem := range *searchres {
		go func(i int, mem namiTypes.SearchMember) {
			id := strconv.Itoa(mem.ID)
			member, err := namigo.GetMemberDetails(id, config.Gruppierung)
			if err != nil {
				log.Println("Error getting details of: " + mem.Vorname + mem.Nachname)
				log.Println(err)
			}
			//1999-01-31 00:00:00
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
			progress <- 1
		}(i, mem)
	}
	wg.Wait()
	return &res
}

func genSepaXml(list *[]namiTypes.Member) {

	fmt.Println("Converting to SEPA format")

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

	helpers.ReadBankInfo()

	for _, member := range *list {
		valid := verifyMandate(&member)
		date, err := time.Parse("02012006", member.Kontoverbindung.Kontonummer)
		if err != nil {
			valid = false
		}
		amount := calcAmount(member.BeitragsartID, member.GenericField1, member.MglType)
		if amount == 0.00 {
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

		} else {
			log.Printf("Mandate verification failed for: %s %s \n", member.Vorname, member.Nachname)

		}

	}

	paym.TransacNb = len(paym.Transactions)
	paym.CtrlSum = f2s(ctrlSum)

	docxml.AddPayment(paym)
	docxml.TransacNb = len(paym.Transactions)
	docxml.CtrlSum = f2s(ctrlSum)

	fmt.Println("Writing SEPA XML file to: ", *outFile)
	file, err := os.Create(*outFile)
	if err != nil {
		log.Fatal("cant create file")
	}
	err = docxml.WriteLatin1(file)
	if err != nil {
		log.Fatal(err)
	}

}

func verifyMandate(member *namiTypes.Member) bool {
	if member.Kontoverbindung.Kontoinhaber == "" {
		return false
	} else if member.Kontoverbindung.Kontonummer == "" {
		return false
	}
	iban, err := iban.NewIBAN(member.Kontoverbindung.Iban)
	if err != nil {
		return false
	}
	if member.Kontoverbindung.Bic == "" {
		if iban.CountryCode != "DE" {
			return false
		}
		code := iban.Code[4:12]
		bic, err := helpers.BicFromCode(code)
		if err != nil {
			return false
		}
		member.Kontoverbindung.Bic = bic
	}
	//All good, cleaning up strings
	member.Kontoverbindung.Kontoinhaber = safeStr(member.Kontoverbindung.Kontoinhaber)
	member.Vorname = safeStr(member.Vorname)
	member.Nachname = safeStr(member.Nachname)

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

var replaceDict = map[string]string{
	"Á": "A", "À": "A", "Â": "A", "Ã": "A", "Å": "A", "Ä": "Ae", "Æ": "AE", "Ç": "C", "É": "E", "È": "E", "Ê": "E", "Ë": "E",
	"Í": "I", "Ì": "I", "Î": "I", "Ï": "I", "Ð": "Eth", "Ñ": "N", "Ó": "O", "Ò": "O", "Ô": "O", "Õ": "O", "Ö": "O", "Ø": "O",
	"Ú": "U", "Ù": "U", "Û": "U", "Ü": "Ue", "Ý": "Y", "á": "a", "à": "a", "â": "a", "ã": "a", "å": "a", "ä": "ae", "æ": "ae",
	"ç": "c", "é": "e", "è": "e", "ê": "e", "ë": "e", "í": "i", "ì": "i", "î": "i", "ï": "i", "ð": "eth", "ñ": "n", "ó": "o",
	"ò": "o", "ô": "o", "õ": "o", "ö": "oe", "ø": "o", "ú": "u", "ù": "u", "û": "u", "ü": "ue", "ý": "y", "ß": "ss", "þ": "thorn",
	"ÿ": "y", "&": "u.", "@": "at", "#": "h", "$": "s", "%": "perc", "^": "-", "*": "-",
}

func safeStr(s string) string {
	for key, val := range replaceDict {
		s = strings.ReplaceAll(s, key, val)
	}
	return s
}

func printProgress(ch *chan int, maxVal int, wg *sync.WaitGroup) {
	defer close(*ch)
	val := 0
	fmt.Printf("Getting member details [%d/%d]", val, maxVal)
	for i := range *ch {
		val += i
		fmt.Print("\033[2K") // Löscht die aktuelle Zeile
		fmt.Printf("\rGetting member details [%d/%d]", val, maxVal)

		if val == maxVal {
			fmt.Printf("\n")
			wg.Done()
		}
	}
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
var cfgFile *string

func readSepaCfg() {
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
	sepaCfg.CreditorName = safeStr(sepaCfg.CreditorName)
	sepaCfg.CollectionInfo = safeStr(sepaCfg.CollectionInfo)
}

func init() {
	rootCmd.AddCommand(sepaCmd)

	occupation = sepaCmd.Flags().StringP("occupation", "o", "", "Occupation (if any) for options see 'occupation' sub command help")
	subdivision = sepaCmd.Flags().StringP("subdivision", "d", "", "Subdivision (if any) for options see 'subdivision' sub command help")
	tag = sepaCmd.Flags().StringP("tag", "t", "", "Tag (if any)")
	searchAll = sepaCmd.Flags().BoolP("all", "a", false, "Create file for ALL members")
	sepaCmd.MarkFlagsOneRequired("occupation", "subdivision", "tag", "all")
	sepaCmd.MarkFlagsMutuallyExclusive("all", "occupation")
	sepaCmd.MarkFlagsMutuallyExclusive("all", "subdivision")
	sepaCmd.MarkFlagsMutuallyExclusive("all", "tag")

	fixedFee = sepaCmd.Flags().Float64("fee", 0.00, "Fixed Fee, ignore member fees and set a fixed fee")

	cfgFile = sepaCmd.Flags().StringP("sepaConfig", "s", "", "Path to the sepa config, default is ~/sepa.yml")

	outFile = sepaCmd.Flags().String("out", "", "Output file")
	sepaCmd.MarkFlagRequired("out")

}
