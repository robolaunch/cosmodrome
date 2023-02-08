package api

import "strconv"

type VDIBase struct {
	BuildComponent
	GPUAgnostic         bool   `yaml:"gpuAgnostic"`
	NvidiaDriverVersion string `yaml:"nvidiaDriverVersion"`
}

func NewVDIBase(gpuAgnostic bool, driverVersion string, ubuntuDistro UbuntuDistro, pushComponent bool) *VDIBase {
	vdiBase := VDIBase{}

	vdiBase.GPUAgnostic = gpuAgnostic
	vdiBase.NvidiaDriverVersion = driverVersion

	vdiBase.BuildComponent.Name = "vdi-base"
	vdiBase.BuildComponent.Image = "driver"
	vdiBase.BuildComponent.Tag = string(ubuntuDistro) + "-" + driverVersion
	vdiBase.BuildComponent.BaseImage = "OPENGL-IMAGE"
	vdiBase.BuildComponent.Directory = "images/vdi/base"
	vdiBase.BuildComponent.DockerfilePath = "images/vdi/base/Dockerfile"
	vdiBase.BuildComponent.BuildArgs = map[string]string{
		"NVIDIA_DRIVER_VERSION": driverVersion,
		"GPU_AGNOSTIC":          strconv.FormatBool(gpuAgnostic),
	}
	vdiBase.BuildComponent.Platforms = []string{"amd64"}
	vdiBase.BuildComponent.PushComponent = pushComponent

	return &vdiBase
}
