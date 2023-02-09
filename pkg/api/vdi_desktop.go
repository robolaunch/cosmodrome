package api

type VDIDesktop struct {
	BuildComponent
	UbuntuDesktop string `yaml:"ubuntuDesktop"`
}

func NewVDIDesktop(ubuntuDesktop string, driverVersion string, ubuntuDistro UbuntuDistro, vdiBaseImage string, pushComponent bool) *VDIDesktop {
	vdiDesktop := VDIDesktop{}

	vdiDesktop.UbuntuDesktop = ubuntuDesktop

	vdiDesktop.BuildComponent.Name = "vdi-desktop"
	vdiDesktop.BuildComponent.Image = "vdi"
	vdiDesktop.BuildComponent.Tag = string(ubuntuDistro) + "-" + driverVersion + "-" + ubuntuDesktop
	vdiDesktop.BuildComponent.BaseImage = vdiBaseImage
	vdiDesktop.BuildComponent.Directory = "images/vdi/" + ubuntuDesktop
	vdiDesktop.BuildComponent.DockerfilePath = "images/vdi/" + ubuntuDesktop + "/Dockerfile"
	vdiDesktop.BuildComponent.BuildArgs = map[string]string{}
	vdiDesktop.BuildComponent.Platforms = []string{"amd64"}
	vdiDesktop.BuildComponent.PushComponent = pushComponent

	return &vdiDesktop
}
