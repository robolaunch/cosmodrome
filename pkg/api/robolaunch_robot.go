package api

type RobotBase struct {
	BuildComponent
	ROSBridgeDistribution1 string `yaml:"bridgeDistro1"`
	ROSBridgeDistribution2 string `yaml:"bridgeDistro2"`
}

func NewRobotBase(rosDistributions []string, ubuntuDesktop string, rosImage string) *RobotBase {
	robotBase := RobotBase{}

	distroStr := ""
	if len(rosDistributions) == 1 {
		distroStr = rosDistributions[0]
		robotBase.ROSBridgeDistribution1 = rosDistributions[0]
		robotBase.ROSBridgeDistribution2 = rosDistributions[0]
	} else if len(rosDistributions) == 1 {
		distroStr = rosDistributions[0] + "-" + rosDistributions[1]
		robotBase.ROSBridgeDistribution1 = rosDistributions[0]
		robotBase.ROSBridgeDistribution2 = rosDistributions[1]
	}

	robotBase.BuildComponent.Name = "robot-base"
	robotBase.BuildComponent.Image = "robot"
	robotBase.BuildComponent.Tag = "base-" + distroStr + "-agnostic-" + ubuntuDesktop
	robotBase.BuildComponent.BaseImage = rosImage
	robotBase.BuildComponent.Directory = "images/robot"
	robotBase.BuildComponent.DockerfilePath = "images/robot/Dockerfile"
	robotBase.BuildComponent.BuildArgs = map[string]string{}
	robotBase.BuildComponent.Platforms = []string{"amd64", "arm64"}

	return &robotBase
}
