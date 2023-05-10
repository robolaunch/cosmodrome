package api

import (
	"errors"
	"reflect"
)

type Image struct {
	Name       string `yaml:"name"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type Step struct {
	Name      string            `yaml:"name"`
	Image     Image             `yaml:"image"`
	BaseStep  string            `yaml:"baseStep"`
	BuildArgs map[string]string `yaml:"buildArgs"`
	Push      bool              `yaml:"push"`
}

func (step *Step) validate() error {
	if reflect.DeepEqual(step.Name, "") {
		return errors.New(".steps[].name cannot be empty")
	}
	if reflect.DeepEqual(step.Image.Repository, "") {
		return errors.New(".steps[].image.repository cannot be empty")
	}
	if reflect.DeepEqual(step.Image.Tag, "") {
		return errors.New(".steps[].image.tag cannot be empty")
	}
	if reflect.DeepEqual(step.Image.Name, "") {
		return errors.New(".steps[].image.name cannot be empty")
	}
	return nil
}

func (step *Step) setImageName(lc LaunchConfig) {
	step.Image.Name = lc.Registry + "/" + lc.Organization + "/" + step.Image.Repository + ":" + step.Image.Tag
}

func (step *Step) start(status *LaunchStatus) error {
	stepStatus := NewStepStatus()
	stepStatus.Step = *step

	// if it's the first step
	if step.BaseStep == "" {
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
