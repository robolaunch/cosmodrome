package api

import "strconv"

func NewVDIBase(ubuntuDistro UbuntuDistro) *BuildComponent {
	vdiBase := BuildComponent{}

	vdiBase.Name = "vdi-base"
	vdiBase.Image = "driver"
	vdiBase.Tag = string(ubuntuDistro) + "-agnostic"
	if ubuntuDistro == UbuntuDistroJammy {
		vdiBase.BaseImage = "robolaunchio/opengl:1.4-runtime-ubuntu22.04"
	} else if ubuntuDistro == UbuntuDistroFocal {
		vdiBase.BaseImage = "nvidia/opengl:1.2-glvnd-runtime-ubuntu20.04"
	}
	vdiBase.Directory = "images/vdi/base"
	vdiBase.DockerfilePath = "images/vdi/base/Dockerfile"
	vdiBase.BuildArgs = map[string]string{
		"NVIDIA_DRIVER_VERSION": "agnostic",
		"GPU_AGNOSTIC":          strconv.FormatBool(true),
	}
	vdiBase.Platforms = []string{"amd64"}

	return &vdiBase
}
