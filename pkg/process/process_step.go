package process

import (
	"context"
	"time"

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
		if err := push(step, stepStatus); err != nil {
			return err
		}
	}

	return nil
}

func build(step *api.Step, baseStep api.Step, stepStatus *api.StepStatus) error {

	stepStatus.Phase = api.StepPhaseBuilding

	// building jobs
	StepLog.Println("Building step: " + step.Name)
	docker.Build(context.Background(), "Dockerfile", "test.local/Dockerfile", baseStep.Image.Name, *step)
	// ***

	return nil
}

func push(step *api.Step, stepStatus *api.StepStatus) error {

	stepStatus.Phase = api.StepPhasePushing

	// pushing jobs
	StepLog.Println("Pushing step: " + step.Name)
	time.Sleep(time.Second * 1)
	// ***

	return nil
}
