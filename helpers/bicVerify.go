package helpers

import (
	"embed"
	"encoding/json"
	"errors"
	"log"
)

type BankInfo struct {
	BankCode     string `json:"bank_code"`
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	Bic          string `json:"bic"`
	Primary      bool   `json:"primary"`
	CountryCode  string `json:"country_code"`
	ChecksumAlgo string `json:"checksum_algo"`
}

//go:embed bics.json
var rawBics embed.FS

var banks []BankInfo

func ReadBankInfo() {
    content, err := rawBics.ReadFile("bics.json")
    if err != nil {
        log.Fatal("Error when opening Bics file: ", err)
    }
 
    // Now let's unmarshall the data into `payload`
    err = json.Unmarshal(content, &banks)
    if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
    }
}

func BicFromCode(bankCode string) (string, error) {
    for _, bank := range banks {
        if bank.BankCode == bankCode {
            return bank.Bic, nil
        }
    }
    return "", errors.New("Bank Code not found!")
    
}
