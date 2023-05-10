package api

import (
	"fmt"
	"time"
)

func (step *Step) start(lc LaunchConfig, status *LaunchStatus) error {
	stepStatus := NewStepStatus()
	stepStatus.Step = *step

	baseStep, err := step.getBaseStep(lc)
	if err != nil {
		return err
	}

	fmt.Println("Starting to process step: " + step.Name)

	if err := step.build(baseStep, stepStatus); err != nil {
		return err
	}
	if step.Push {
		if err := step.push(stepStatus); err != nil {
			return err
		}
	}

	return nil
}

func (step *Step) build(baseStep Step, stepStatus *StepStatus) error {

	stepStatus.Phase = StepPhaseBuilding

	// building jobs
	fmt.Println("Building step: " + step.Name)
	time.Sleep(time.Second * 1)
	// ***

	return nil
}

func (step *Step) push(stepStatus *StepStatus) error {

	stepStatus.Phase = StepPhasePushing

	// pushing jobs
	fmt.Println("Pushing step: " + step.Name)
	time.Sleep(time.Second * 1)
	// ***

	return nil
}
