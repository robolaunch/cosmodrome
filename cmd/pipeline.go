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

	var rosDistributions []api.ROSDistro

	if multipleROSDistro {

		rosDistro1, err := askCustomSelectable("First ROS Distro", getMultipleDistroList(nil))
		if err != nil {
			return nil, err
		}

		rosDistro2, err := askCustomSelectable("Second ROS Distro", getMultipleDistroList((*api.ROSDistro)(&rosDistro1)))
		if err != nil {
			return nil, err
		}

		rosDistributions = []api.ROSDistro{api.ROSDistro(rosDistro1), api.ROSDistro(rosDistro2)}

	} else {

		rosDistro, err := askCustomSelectable("ROS Distro", []string{string(api.ROSDistroHumble), string(api.ROSDistroFoxy), string(api.ROSDistroGalactic)})
		if err != nil {
			return nil, err
		}

		rosDistributions = []api.ROSDistro{api.ROSDistro(rosDistro)}
	}

	var ubuntuDistro api.UbuntuDistro
	for _, v := range rosDistributions {
		switch v {
		case api.ROSDistroHumble:
			ubuntuDistro = api.UbuntuDistroJammy
		case api.ROSDistroFoxy:
			ubuntuDistro = api.UbuntuDistroFocal
		case api.ROSDistroGalactic:
			ubuntuDistro = api.UbuntuDistroFocal
		}
	}

	ubuntuDesktop, err := askCustomSelectable("Ubuntu Desktop", []string{"xfce", "mate"})
	if err != nil {
		return nil, err
	}

	pipeline := api.NewPipeline(name, registry, rosDistributions, ubuntuDistro, ubuntuDesktop)

	return pipeline, nil
}

func getMultipleDistroList(distro *api.ROSDistro) []string {
	if distro == nil {
		return []string{string(api.ROSDistroFoxy), string(api.ROSDistroGalactic)}
	} else {
		switch *distro {
		case api.ROSDistroFoxy:
			return []string{string(api.ROSDistroGalactic)}
		case api.ROSDistroGalactic:
			return []string{string(api.ROSDistroFoxy)}
		}
	}

	return []string{}
}
