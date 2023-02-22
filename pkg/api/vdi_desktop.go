package api

func NewVDIDesktop(ubuntuDesktop string, ubuntuDistro UbuntuDistro, vdiBaseImage string) *BuildComponent {
	vdiDesktop := BuildComponent{}

	vdiDesktop.Name = "vdi-desktop"
	vdiDesktop.Image = "vdi"
	vdiDesktop.Tag = string(ubuntuDistro) + "-agnostic-" + ubuntuDesktop
	vdiDesktop.BaseImage = vdiBaseImage
	vdiDesktop.Directory = "images/vdi/" + ubuntuDesktop
	vdiDesktop.DockerfilePath = "images/vdi/" + ubuntuDesktop + "/Dockerfile"
	vdiDesktop.BuildArgs = map[string]string{}
	vdiDesktop.Platforms = []string{"amd64"}

	return &vdiDesktop
}
