package process

import (
	"context"

	"github.com/robolaunch/cosmodrome/pkg/api"
	"github.com/robolaunch/cosmodrome/pkg/docker"
)

func start(step *api.Step, status *api.StepStatus, lc api.LaunchConfig) error {

	status.Step = *step

	baseStep, err := step.GetBaseStep(lc)
	if err != nil {
		return err
	}

	if err := build(step, baseStep, status); err != nil {
		status.Phase = api.StepPhaseFailed
		return err
	}
	if step.Push {
		if err := push(step, status, lc); err != nil {
			status.Phase = api.StepPhaseFailed
			return err
		}
	}

	status.Phase = api.StepPhaseSucceeded

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
