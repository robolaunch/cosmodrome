package process

import (
	"context"

	"github.com/robolaunch/cosmodrome/pkg/api"
	"github.com/robolaunch/cosmodrome/pkg/docker"
)

func start(step *api.Step, lc api.LaunchConfig, status *api.LaunchStatus) error {

	stepStatus := api.NewStepStatus()
	stepStatus.Step = *step

	baseStep, err := step.GetBaseStep(lc)
	if err != nil {
		return err
	}

	if err := build(step, baseStep, stepStatus); err != nil {
		return err
	}
	if step.Push {
		if err := push(step, stepStatus, lc); err != nil {
			return err
		}
	}

	return nil
}

func build(step *api.Step, baseStep api.Step, stepStatus *api.StepStatus) error {

	stepStatus.Phase = api.StepPhaseBuilding

	StepLog.Println("Building step: " + step.Name)
	if err := docker.Build(context.Background(), "Dockerfile", step.Dockerfile, baseStep.Image.Name, *step); err != nil {
		return err
	}

	return nil
}

func push(step *api.Step, stepStatus *api.StepStatus, lc api.LaunchConfig) error {

	stepStatus.Phase = api.StepPhasePushing

	StepLog.Println("Pushing step: " + step.Name)
	if err := docker.Push(context.Background(), *step, lc); err != nil {
		return err
	}

	return nil
}
