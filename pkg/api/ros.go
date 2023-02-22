package api

func NewROS(rosDistributions []ROSDistro, ubuntuDesktop string, vdiDesktopImage string) *BuildComponent {
	ros := BuildComponent{}

	distroStr := ""
	if len(rosDistributions) == 1 {
		distroStr = string(string(rosDistributions[0]))
	} else if len(rosDistributions) == 1 {
		distroStr = string(string(rosDistributions[0])) + "-" + string(string(rosDistributions[1]))
	}

	ros.Name = "ros"
	ros.Image = distroStr
	ros.Tag = "agnostic-" + ubuntuDesktop
	ros.BaseImage = vdiDesktopImage
	ros.Directory = "images/ros/" + distroStr
	ros.DockerfilePath = "images/ros/" + distroStr + "/Dockerfile"
	ros.BuildArgs = map[string]string{}
	ros.Platforms = []string{"amd64", "arm64"}

	return &ros
}
