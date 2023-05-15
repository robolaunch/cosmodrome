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
	Logfile      string `yaml:"logfile,omitempty"`
	Verbose      bool   `yaml:"verbose,omitempty"`
	NoCache      bool   `yaml:"nocache,omitempty"`
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
		lc.Steps[k].Default(*lc)
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
	for k, step := range lc.Steps {
		if err := step.validate(k); err != nil {
			return err
		}
	}

	if err := lc.validateStepsSemantics(); err != nil {
		return err
	}
	return nil
}

func (lc *LaunchConfig) validateStepsSemantics() error {
	// check if base step is disabled for first step
	if !reflect.DeepEqual(lc.Steps[0].BaseStep, "") {
		return errors.New(".steps[0].baseStep should be empty")
	}

	// check if base step is populated
	for i, step := range lc.Steps[1:] {
		if reflect.DeepEqual(step.BaseStep, "") {
			return errors.New(".steps[" + strconv.Itoa(i+1) + "].baseStep cannot be empty")
		}
	}

	// check if base step really exists
	for i := range lc.Steps {
		if i == 0 {
			continue
		} else {
			hasBaseStep := false
			for j := 0; j < i; j++ {
				if lc.Steps[i].BaseStep == lc.Steps[j].Name {
					hasBaseStep = true
					break
				}
			}

			if !hasBaseStep {
				return errors.New(".steps[" + strconv.Itoa(i) + "].baseStep is invalid. no step named `" + lc.Steps[i].BaseStep + "` is found before the step `" + lc.Steps[i].Name + "`.")
			}
		}
	}

	return nil
}
