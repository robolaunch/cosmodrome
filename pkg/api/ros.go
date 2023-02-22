package api

type ROS struct {
	BuildComponent
	ROSDistributions []string `yaml:"rosDistributions"`
}

func NewROS(rosDistributions []string, ubuntuDesktop string, vdiDesktopImage string) *ROS {
	ros := ROS{}

	distroStr := ""
	if len(rosDistributions) == 1 {
		distroStr = rosDistributions[0]
	} else if len(rosDistributions) == 1 {
		distroStr = rosDistributions[0] + "-" + rosDistributions[1]
	}

	ros.ROSDistributions = rosDistributions

	ros.BuildComponent.Name = "ros"
	ros.BuildComponent.Image = distroStr
	ros.BuildComponent.Tag = "agnostic-" + ubuntuDesktop
	ros.BuildComponent.BaseImage = vdiDesktopImage
	ros.BuildComponent.Directory = "images/ros/" + distroStr
	ros.BuildComponent.DockerfilePath = "images/ros/" + distroStr + "/Dockerfile"
	ros.BuildComponent.BuildArgs = map[string]string{}
	ros.BuildComponent.Platforms = []string{"amd64", "arm64"}

	return &ros
}
