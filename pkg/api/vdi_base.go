package api

import "strconv"

type VDIBase struct {
	BuildComponent
}

func NewVDIBase(ubuntuDistro UbuntuDistro) *VDIBase {
	vdiBase := VDIBase{}

	vdiBase.BuildComponent.Name = "vdi-base"
	vdiBase.BuildComponent.Image = "driver"
	vdiBase.BuildComponent.Tag = string(ubuntuDistro) + "-agnostic"
	vdiBase.BuildComponent.BaseImage = "OPENGL-IMAGE"
	vdiBase.BuildComponent.Directory = "images/vdi/base"
	vdiBase.BuildComponent.DockerfilePath = "images/vdi/base/Dockerfile"
	vdiBase.BuildComponent.BuildArgs = map[string]string{
		"NVIDIA_DRIVER_VERSION": "agnostic",
		"GPU_AGNOSTIC":          strconv.FormatBool(true),
	}
	vdiBase.BuildComponent.Platforms = []string{"amd64"}

	return &vdiBase
}
