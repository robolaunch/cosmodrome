/*
Copyright © 2023 robolaunch
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/robolaunch/cosmodrome/pkg/api"
	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("pipeline called")

		export, _ := cmd.Flags().GetBool("export")
		view, _ := cmd.Flags().GetBool("view")
		if export && view {
			panic(errors.New("you cannot use export flag and view flag at the same time"))
		}

		name, err := askName()
		if err != nil {
			panic(err)
		}

		registry, err := askRegistry()
		if err != nil {
			panic(err)
		}

		pushComponents, err := askPushComponents()
		if err != nil {
			panic(err)
		}

		pipeline := api.NewPipeline(name, registry, pushComponents)

		if view {

			err = pipeline.View()
			if err != nil {
				panic(err)
			}

		} else {

			err = pipeline.Export()
			if err != nil {
				panic(err)
			}

		}
	},
}

func init() {
	createCmd.AddCommand(pipelineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	pipelineCmd.Flags().BoolP("export", "e", false, "Export pipeline configuration to a YAML file.")
	pipelineCmd.Flags().BoolP("view", "v", false, "View pipeline configuration.")
}

func askName() (string, error) {
	var name string

	fmt.Print("Enter pipeline name (eg. sunday): ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		return "", err
	}

	// validate

	return name, nil
}

func askRegistry() (string, error) {
	var registry string

	fmt.Print("Enter registry (eg. robolaunchio): ")
	_, err := fmt.Scanln(&registry)
	if err != nil {
		return "", err
	}

	// validate

	return registry, nil
}

func askPushComponents() (bool, error) {
	var pushComponents string

	fmt.Print("Push components to the registry (y/n): ")
	_, err := fmt.Scanln(&pushComponents)
	if err != nil {
		return false, err
	}

	// validate

	// convert

	switch pushComponents {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, errors.New("wrong format")
	}

}
