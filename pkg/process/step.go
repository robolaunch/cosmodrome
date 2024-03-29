package process

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/robolaunch/cosmodrome/pkg/api"
	"github.com/robolaunch/cosmodrome/pkg/docker"
)

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	go func() {
		for sig := range c {
			fmt.Println("Stopped process: " + sig.String())
			terminateSpinner()
			os.Exit(1)
		}
	}()
}

func terminateSpinner() {
	if logSpinner.Enabled() {
		logSpinner.Stop()
	}
}

func start(step *api.Step, status *api.StepStatus, lc api.LaunchConfig) error {

	go handleSignal()
	status.Step = *step

	baseStep, err := step.GetBaseStep(lc)
	if err != nil {
		terminateSpinner()
		return err
	}

	if err := build(step, baseStep, status, lc); err != nil {
		status.Phase = api.StepPhaseFailed
		terminateSpinner()
		return err
	}
	if step.Push && len(step.Platforms) == 0 {
		if err := push(step, status, lc); err != nil {
			status.Phase = api.StepPhaseFailed
			terminateSpinner()
			return err
		}
	}

	status.Phase = api.StepPhaseSucceeded

	return nil
}

func build(step *api.Step, baseStep api.Step, stepStatus *api.StepStatus, lc api.LaunchConfig) error {

	stepStatus.Phase = api.StepPhaseBuilding
	if len(step.Platforms) == 0 {
		GetSpinner(StepLog, " Building step: "+step.Name)
		logSpinner.Start()
		if err := docker.Build(context.Background(), step.Dockerfile, step.Path, step.Context, *step, lc); err != nil {
			return err
		}
	} else {
		GetSpinner(StepLog, " Building step: "+step.Name)
		if step.Push {
			GetSpinner(StepLog, " Building and pushing step: "+step.Name)
		}
		logSpinner.Start()
		if err := docker.BuildMultiplatform(context.Background(), step.Dockerfile, step.Path, step.Context, baseStep.Image.Name, *step, lc); err != nil {
			return err
		}
	}
	logSpinner.Stop()

	return nil
}

func push(step *api.Step, stepStatus *api.StepStatus, lc api.LaunchConfig) error {

	stepStatus.Phase = api.StepPhasePushing

	GetSpinner(StepLog, " Pushing step: "+step.Name)
	logSpinner.Start()
	if err := docker.Push(context.Background(), *step, lc); err != nil {
		return err
	}
	logSpinner.Stop()

	return nil
}
