/*
Copyright © 2023 robolaunch
*/
package cmd

import (
	"errors"

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

		pipeline, err := askPipelineConfig()
		if err != nil {
			panic(err)
		}

		buildVDI, err := askBinaryQuestion("Build VDI")
		if err != nil {
			panic(err)
		}

		vdiBase := &api.BuildComponent{}
		vdiDesktop := &api.BuildComponent{}

		if buildVDI {

			vdiBase = api.NewVDIBase(pipeline.UbuntuDistro)
			vdiDesktop = api.NewVDIDesktop(pipeline.UbuntuDesktop, pipeline.UbuntuDistro, vdiBase.GetImage(pipeline.Registry))

			pipeline.Components = append(pipeline.Components, *vdiBase)
			pipeline.Components = append(pipeline.Components, *vdiDesktop)

		}

		ros := api.NewROS(pipeline.ROSDistributions, pipeline.UbuntuDesktop, vdiDesktop.GetImage(pipeline.Registry))
		pipeline.Components = append(pipeline.Components, *ros)

		robotBase := api.NewRobotBase(pipeline.ROSDistributions, pipeline.UbuntuDesktop, ros.GetImage(pipeline.Registry))
		pipeline.Components = append(pipeline.Components, *robotBase)

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

	pipelineCmd.Flags().BoolP("export", "e", false, "Export pipeline configuration to a YAML file.")
	pipelineCmd.Flags().BoolP("view", "v", false, "View pipeline configuration.")
}

func askPipelineConfig() (*api.Pipeline, error) {
	name, err := askStringQuestion("Enter pipeline name (eg. sunday)")
	if err != nil {
		panic(err)
	}

	registry, err := askStringQuestion("Enter registry (eg. robolaunchio)")
	if err != nil {
		panic(err)
	}

	multipleROSDistro, err := askBinaryQuestion("Multiple ROS Distro")
	if err != nil {
		panic(err)
	}

	var rosDistro string

	if multipleROSDistro {
		// ask two distro
	} else {
		rosDistro, err = askCustomSelectable("ROS Distro", []string{"humble", "foxy", "galactic"})
		if err != nil {
			return nil, err
		}

	}

	ubuntuDistro, err := askCustomSelectable("Ubuntu Distro", []string{"jammy", "focal"})
	if err != nil {
		panic(err)
	}

	ubuntuDesktop, err := askCustomSelectable("Ubuntu Desktop", []string{"xfce", "mate"})
	if err != nil {
		return nil, err
	}

	pipeline := api.NewPipeline(name, registry, []api.ROSDistro{api.ROSDistro(rosDistro)}, api.UbuntuDistro(ubuntuDistro), ubuntuDesktop)

	return pipeline, nil
}
