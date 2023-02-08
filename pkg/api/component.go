package api

import (
	"fmt"
)

type BuildComponentInterface interface {
	Build()
	Push()
	View()
}

type BuildComponent struct {
	BuildComponentInterface `yaml:",omitempty"`
	Name                    string            `yaml:"name"`
	Image                   string            `yaml:"image"`
	Tag                     string            `yaml:"tag"`
	BaseImage               string            `yaml:"baseImage"`
	Directory               string            `yaml:"directory"`
	DockerfilePath          string            `yaml:"dockerfilePath"`
	BuildArgs               map[string]string `yaml:"buildArgs"`
	Platforms               []string          `yaml:"platforms"`
	PushComponent           bool              `yaml:"pushComponent"`
}

func (bc *BuildComponent) Build() {
	fmt.Println(bc.Image, "Component is built.")
}

func (bc *BuildComponent) Push() {
	fmt.Println(bc.Image, "Component is pushed.")
}

func (bc *BuildComponent) View() {
	fmt.Println(bc.Image, "Component is viewed.")
}
