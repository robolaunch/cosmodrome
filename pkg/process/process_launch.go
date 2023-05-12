package process

import "github.com/robolaunch/cosmodrome/pkg/api"

func Start(lc *api.LaunchConfig) error {
	status := api.NewLaunchStatus()

	for k := range lc.Steps {
		api.LaunchLog.Println("Starting to process step: " + lc.Steps[k].Name)
		if err := start(&lc.Steps[k], *lc, status); err != nil {
			return err
		}
	}

	return nil
}
