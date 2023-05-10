package api

import (
	"errors"
	"reflect"
)

type Step struct {
	Name      string            `yaml:"name"`
	RootImage string            `yaml:"rootImage"`
	BaseStep  string            `yaml:"baseStep"`
	BuildArgs map[string]string `yaml:"buildArgs"`
	Push      bool              `yaml:"push"`
}

func (step *Step) validate() error {
	if reflect.DeepEqual(step.Name, "") {
		return errors.New(".steps[].name cannot be empty")
	}
	return nil
}
