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

type PipelineAbstract struct {
	PipelineInterface `yaml:",omitempty"`
	Name              string                    `yaml:"name"`
	Registry          string                    `yaml:"registry"`
	UbuntuDistro      UbuntuDistro              `yaml:"ubuntuDistro"`
	Components        []BuildComponentInterface `yaml:"components,omitempty"`
}

type Pipeline struct {
	PipelineAbstract `yaml:"pipeline,omitempty"`
}

func NewPipeline(name string, registry string, ubuntuDistro UbuntuDistro) *Pipeline {
	pipeline := Pipeline{}

	pipeline.Name = name
	pipeline.Registry = registry
	pipeline.UbuntuDistro = ubuntuDistro
	pipeline.Components = []BuildComponentInterface{}

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
