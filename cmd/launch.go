/*
Copyright © 2023 robolaunch
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/robolaunch/cosmodrome/pkg/api"
	"github.com/robolaunch/cosmodrome/pkg/process"
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

		// check if config will be printed out
		printConfig, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			Error.Println(err.Error())
			os.Exit(2)
		}

		// get launch config and convert it to struct
		launchCfg := &api.LaunchConfig{}
		err = viper.Unmarshal(launchCfg)
		if err != nil {
			fmt.Printf("unable to decode into config struct, %v", err)
		}

		// default fields
		launchCfg.Default()

		// validate launch config
		if err := launchCfg.Validate(); err != nil {
			Error.Println(err.Error())
			os.Exit(2)
		}

		// print launch config
		if printConfig {
			launchCfg.PrintYAML()
		}

		// process launch
		if err := process.Start(launchCfg); err != nil {
			Error.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	launchCmd.PersistentFlags().BoolP("no-cache", "nc", false, "use `--no-cache` flag when building an image")
	rootCmd.AddCommand(launchCmd)
}
