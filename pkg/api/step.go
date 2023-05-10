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

func (step *Step) start(status *LaunchStatus) error {
	stepStatus := NewStepStatus()
	stepStatus.Step = *step

	// if it's the first step
	if step.RootImage != "" {
		// process first step
	} else {
		// process step
	}

	return nil
}

func (step *Step) build(baseImage, stepStatus *StepStatus) error {

	stepStatus.Phase = StepPhaseBuilding

	// building jobs
	// ***

	return nil
}

func (step *Step) push(stepStatus *StepStatus) error {

	stepStatus.Phase = StepPhasePushing

	// pushing jobs
	// ***

	return nil
}
