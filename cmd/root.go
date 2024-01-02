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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/thisni1s/nami-go"
	namiTypes "github.com/thisni1s/nami-go/types"
	"gopkg.in/yaml.v3"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nami-cli",
	Short: "Command Line Application for interacting with the DPSG NaMi",
	Long: `nami-cli allows you to interact with the DPSG NaMi and get information about your Members.
A config file containing your credentials is needed for this to work.
By default the file should be located at ".nami.yml" and look like this:

username: 133337 # your nami id
password: verysecure
gruppierung: 010101 # your "Stammesnummer"
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var config namiTypes.Config 
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .nami.yaml)")
	if cfgFile == "" {
		cfgFile = ".nami.yml"
	}

	file, err := os.ReadFile(cfgFile)
	if err != nil {
        log.Println("Failed to read config file!")
        log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
        log.Println("Failed to read config file!")
        log.Fatal(err)
	}


	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Login() {
    err := namigo.Login(config.Username, config.Password)
    if err != nil {
        log.Println("Failed to login!")
        log.Fatal(err)
    }
}

func GetGroupId() string {
    return config.Gruppierung
}
