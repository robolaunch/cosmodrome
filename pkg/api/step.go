package api

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Image struct {
	Name       string `yaml:"name"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type Step struct {
	Name       string             `yaml:"name"`
	Image      Image              `yaml:"image"`
	Path       string             `yaml:"path"`
	Dockerfile string             `yaml:"dockerfile"`
	BaseStep   string             `yaml:"baseStep"`
	BuildArgs  map[string]*string `yaml:"buildArgs"`
	Context    string             `yaml:"context"`
	Platforms  []string           `yaml:"platforms"`
	Push       bool               `yaml:"push"`
}

func (step *Step) Default(lc LaunchConfig) {
	step.setImageName(lc)
	step.setContext(lc)
	step.setBuildArgs(lc)
	step.setBaseImage(lc)
}

func (step *Step) setImageName(lc LaunchConfig) {
	step.Image.Name = lc.Registry + "/" + lc.Organization + "/" + step.Image.Repository + ":" + FormatTag(step.Image.Tag) + "-" + lc.Version
}

func FormatTag(tag string) string {
	// no "+" is allowed
	// no "~" is allowed
	tag = strings.ReplaceAll(tag, "+", "-")
	tag = strings.ReplaceAll(tag, "~", "-")
	tag = strings.ReplaceAll(tag, ":", "-")
	return tag
}

func (step *Step) setContext(lc LaunchConfig) {
	if step.Context == "" {
		step.Context = "."
	}
}

func (step *Step) setBuildArgs(lc LaunchConfig) {
	if len(step.BuildArgs) == 0 {
		step.BuildArgs = make(map[string]*string)
	}
}

func (step *Step) setBaseImage(lc LaunchConfig) {
	if _, ok := step.BuildArgs["base_image"]; !ok {
		baseStep, _ := step.GetBaseStep(lc)
		step.BuildArgs["base_image"] = &baseStep.Image.Name
	}
}

func (step *Step) validate(key int) error {
	if reflect.DeepEqual(step.Name, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].name cannot be empty")
	}
	if reflect.DeepEqual(step.Path, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].path cannot be empty")
	}
	if reflect.DeepEqual(step.Dockerfile, "") {
		return errors.New(".steps[" + strconv.Itoa(key) + "].dockerfile cannot be empty")
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

func (step *Step) GetBaseStep(lc LaunchConfig) (Step, error) {

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
