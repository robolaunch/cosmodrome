package api

func NewRobotBase(rosDistributions []ROSDistro, ubuntuDesktop string, rosImage string) *BuildComponent {
	robotBase := BuildComponent{}

	distroStr := ""
	if len(rosDistributions) == 1 {
		distroStr = string(rosDistributions[0])
		robotBase.BuildArgs = map[string]string{
			"BRIDGE_DISTRO_1": string(rosDistributions[0]),
			"BRIDGE_DISTRO_2": string(rosDistributions[0]),
		}
	} else if len(rosDistributions) == 2 {
		distroStr = string(rosDistributions[0]) + "-" + string(rosDistributions[1])
		robotBase.BuildArgs = map[string]string{
			"BRIDGE_DISTRO_1": string(rosDistributions[0]),
			"BRIDGE_DISTRO_2": string(rosDistributions[1]),
		}
	}

	robotBase.Name = "robot-base"
	robotBase.Image = "robot"
	robotBase.Tag = "base-" + distroStr + "-agnostic-" + ubuntuDesktop
	robotBase.BaseImage = rosImage
	robotBase.Directory = "images/robot"
	robotBase.DockerfilePath = "images/robot/Dockerfile"
	robotBase.Platforms = []string{"amd64", "arm64"}

	return &robotBase
}
