package api

import (
	"fmt"
)

type BuildComponentInterface interface {
	Build()
	Push()
	View()
	GetImage(registry string) string
}

type BuildComponent struct {
	BuildComponentInterface
	Name           string            `mapstructure:"name"`
	Image          string            `mapstructure:"image"`
	Tag            string            `mapstructure:"tag"`
	BaseImage      string            `mapstructure:"baseImage"`
	Directory      string            `mapstructure:"directory"`
	DockerfilePath string            `mapstructure:"dockerfilePath"`
	BuildArgs      map[string]string `mapstructure:"buildArgs"`
	Platforms      []string          `mapstructure:"platforms"`
}

func (bc *BuildComponent) Build() {

	execBashCmd("")

}

func (bc *BuildComponent) Push() {
	fmt.Println(bc.Name, "Component is pushed.")
}

func (bc *BuildComponent) View() {
	fmt.Println(bc.Name, "Component is viewed.")
}

func (bc *BuildComponent) GetImage(registry string) string {
	return registry + "/" + bc.Image + ":" + bc.Tag
}
