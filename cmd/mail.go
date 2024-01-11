package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// mailCmd represents the mail command
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Send E-Mails to different Members",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mail called")
	},
}

var MailTag *string

func init() {
	rootCmd.AddCommand(mailCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mailCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	MailTag = mailCmd.PersistentFlags().StringP("tag", "t", "", "Send E-Mail only to users with this tag.")
}

type SepaMailbox struct {
	Name      string
	ChildName string
	Address   string
	Amount    string
	Reference string
}
