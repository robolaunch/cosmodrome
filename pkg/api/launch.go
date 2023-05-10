package api

import (
	"errors"
	"reflect"
)

type LaunchConfig struct {
	Name  string `yaml:"name,omitempty"`
	Steps []Step `yaml:"steps,omitempty"`
}

type Step struct {
	Name      string            `yaml:"name,omitempty"`
	Root      bool              `yaml:"root,omitempty"`
	BaseStep  string            `yaml:"baseStep,omitempty"`
	BuildArgs map[string]string `yaml:"buildArgs,omitempty"`
	Push      bool              `yaml:"push,omitempty"`
}

func (lc *LaunchConfig) Validate() error {
	if reflect.DeepEqual(len(lc.Steps), 0) {
		return errors.New("launch config should contain at least one step")
	}
	return nil
}

func (step *Step) Validate() error {
	if reflect.DeepEqual(step.Name, "") {
		return errors.New(".steps[].name cannot be empty")
	}
	return nil
}
