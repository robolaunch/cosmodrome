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
	BuildComponentInterface
	Image          string
	BaseImage      string
	Directory      string
	DockerfilePath string
	BuildArgs      []string
	Platforms      []string
	Rebuild        bool
	PushComponent  bool
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
