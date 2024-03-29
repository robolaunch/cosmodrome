package process

import (
	"fmt"

	"github.com/kyokomi/emoji/v2"
	"github.com/robolaunch/cosmodrome/pkg/api"
)

func Start(lc *api.LaunchConfig) error {
	status := api.NewLaunchStatus()

	for k := range lc.Steps {
		LaunchLog.Println("Processing step: " + lc.Steps[k].Name)
		stepStatus := api.NewStepStatus()
		if err := start(&lc.Steps[k], stepStatus, *lc); err != nil {
			status.StepStatuses = append(status.StepStatuses, *stepStatus)
			return err
		}

		status.StepStatuses = append(status.StepStatuses, *stepStatus)

		emoji.Print(SuccessLog.Sprint(":whale: " + lc.Steps[k].Image.Name + " is generated"))
		if lc.Steps[k].Push {
			SuccessLog.Print(" and pushed")
		}
		SuccessLog.Println(".")
	}

	fmt.Println()
	emoji.Println(SuccessLog.Sprint(":rocket: Steps are completed."))

	return nil
}
