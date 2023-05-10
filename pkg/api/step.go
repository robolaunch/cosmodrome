package api

import (
	"errors"
	"reflect"
	"strconv"
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

func (step *Step) Default(lc LaunchConfig) {
	step.setImageName(lc)
}

func (step *Step) setImageName(lc LaunchConfig) {
	step.Image.Name = lc.Registry + "/" + lc.Organization + "/" + step.Image.Repository + ":" + step.Image.Tag
}

func (step *Step) validate(key int) error {
	if reflect.DeepEqual(step.Name, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].name cannot be empty")
	}
	if reflect.DeepEqual(step.Image.Repository, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].image.repository cannot be empty")
	}
	if reflect.DeepEqual(step.Image.Tag, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].image.tag cannot be empty")
	}
	if reflect.DeepEqual(step.Image.Name, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].image.name cannot be empty")
	}
	return nil
}

func (step *Step) getBaseStep(lc LaunchConfig) (Step, error) {

	if reflect.DeepEqual(step.BaseStep, "") {
		return Step{}, nil
	}

	for _, s := range lc.Steps {
		if s.Name == step.BaseStep {
			return s, nil
		}
	}

	return Step{}, errors.New("cannot find base step of " + step.Name)
}
