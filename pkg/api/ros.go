package api

type ROS struct {
	BuildComponent
	ROSDistributions []string `yaml:"rosDistributions"`
}
