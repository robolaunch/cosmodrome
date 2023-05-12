package process

import "github.com/robolaunch/cosmodrome/pkg/api"

func Start(lc *api.LaunchConfig) error {
	status := api.NewLaunchStatus()

	for k := range lc.Steps {
		LaunchLog.Println("Starting to process step: " + lc.Steps[k].Name)
		stepStatus := api.NewStepStatus()
		if err := start(&lc.Steps[k], stepStatus, *lc); err != nil {
			status.StepStatuses = append(status.StepStatuses, *stepStatus)
			return err
		}

		status.StepStatuses = append(status.StepStatuses, *stepStatus)
	}

	return nil
}
