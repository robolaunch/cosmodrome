package api

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

type LaunchConfig struct {
	Name         string `yaml:"name,omitempty"`
	Registry     string `yaml:"registry,omitempty"`
	Organization string `yaml:"organization,omitempty"`
	Steps        []Step `yaml:"steps,omitempty"`
}

func (lc *LaunchConfig) PrintYAML() error {
	yamlData, err := yaml.Marshal(lc)
	if err != nil {
		return err
	}
	fmt.Println(string(yamlData))
	return nil
}

func (lc *LaunchConfig) Default() {
	for k := range lc.Steps {
		lc.Steps[k].setImageName(*lc)
	}
}

func (lc *LaunchConfig) Validate() error {
	if reflect.DeepEqual(len(lc.Steps), 0) {
		return errors.New(".steps should contain at least one step")
	}
	if reflect.DeepEqual(lc.Registry, "") {
		return errors.New(".registry cannot be empty")
	}
	if reflect.DeepEqual(lc.Organization, "") {
		return errors.New(".organization cannot be empty")
	}
	if err := lc.validateSteps(); err != nil {
		return err
	}
	return nil
}

func (lc *LaunchConfig) validateSteps() error {
	for _, step := range lc.Steps {
		if err := step.validate(); err != nil {
			return err
		}
	}

	if err := lc.validateStepsSemantics(); err != nil {
		return err
	}
	return nil
}

func (lc *LaunchConfig) validateStepsSemantics() error {
	if !reflect.DeepEqual(lc.Steps[0].BaseStep, "") {
		return errors.New(".steps[0].baseStep should be empty")
	}

	for i, step := range lc.Steps[1:] {
		if reflect.DeepEqual(step.BaseStep, "") {
			return errors.New(".steps[" + strconv.Itoa(i+1) + "].baseStep cannot be empty")
		}
	}
	return nil
}
