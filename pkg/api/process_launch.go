package api

func (lc *LaunchConfig) Start() error {
	status := NewLaunchStatus()

	for k := range lc.Steps {
		if err := lc.Steps[k].start(*lc, status); err != nil {
			return err
		}
	}

	return nil
}
