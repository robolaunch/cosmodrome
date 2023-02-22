package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type PipelineInterface interface {
	Build()
	Push()
	View()
	Export()
}

type UbuntuDistro string

const (
	UbuntuDistroFocal UbuntuDistro = "focal"
	UbuntuDistroJammy UbuntuDistro = "jammy"
)

type ROSDistro string

const (
	ROSDistroHumble   ROSDistro = "humble"
	ROSDistroFoxy     ROSDistro = "foxy"
	ROSDistroGalactic ROSDistro = "galactic"
)

// type Components struct {
// 	VDIBase    VDIBase    `mapstructure:"vdiBase,omitempty"`
// 	VDIDesktop VDIDesktop `mapstructure:"vdiDesktop,omitempty"`
// 	ROS        ROS        `mapstructure:"ros,omitempty"`
// 	RobotBase  RobotBase  `mapstructure:"robotBase,omitempty"`
// }

type Pipeline struct {
	PipelineInterface `json:"-"`
	Name              string           `mapstructure:"name"`
	Registry          string           `mapstructure:"registry"`
	UbuntuDistro      UbuntuDistro     `mapstructure:"ubuntuDistro"`
	ROSDistributions  []ROSDistro      `mapstructure:"rosDistributions"`
	UbuntuDesktop     string           `mapstructure:"ubuntuDesktop"`
	Components        []BuildComponent `mapstructure:"components"`
}

func NewPipeline(name string, registry string, rosDistributions []ROSDistro, ubuntuDistro UbuntuDistro, ubuntuDesktop string) *Pipeline {
	pipeline := Pipeline{}

	pipeline.Name = name
	pipeline.Registry = registry
	pipeline.UbuntuDistro = ubuntuDistro
	pipeline.ROSDistributions = rosDistributions
	pipeline.UbuntuDesktop = ubuntuDesktop
	pipeline.Components = []BuildComponent{}

	return &pipeline
}

func (p *Pipeline) Build() {
	fmt.Println(p.Name, "Pipeline is built.")
}

func (p *Pipeline) Push() {
	fmt.Println(p.Name, "Pipeline is pushed.")
}

func (p *Pipeline) View() error {
	pipelineYAML, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}

	fmt.Println("\n" + string(pipelineYAML))
	return nil
}

func (p *Pipeline) Export() error {

	pipelineYAML, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}

	directory := "pipelines"
	_ = os.Mkdir(directory, os.ModePerm)

	fileName := filepath.Join(directory, p.Name+".yaml")
	err = ioutil.WriteFile(fileName, pipelineYAML, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Pipeline is exported to", fileName+".")

	return nil
}
