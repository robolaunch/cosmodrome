/*
Copyright © 2023 robolaunch
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/robolaunch/cosmodrome/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch verb starts building components.",
	Long:  `A longer description for launch.`,
	Run: func(cmd *cobra.Command, args []string) {

		// get launch config and convert it to struct
		launchCfg := &api.LaunchConfig{}
		err := viper.Unmarshal(launchCfg)
		if err != nil {
			fmt.Printf("unable to decode into config struct, %v", err)
		}

		// validate launch config
		if err := launchCfg.Validate(); err != nil {
			Error.Println(err.Error())
			os.Exit(2)
		}

	},
}

func init() {
	rootCmd.AddCommand(launchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// launchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// launchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
