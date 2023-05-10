package api

type StepPhase string

const (
	StepPhaseNotStarted StepPhase = "NotStarted"
	StepPhasePulling    StepPhase = "Pulling"
	StepPhaseBuilding   StepPhase = "Building"
	StepPhasePushing    StepPhase = "Pushing"
	StepPhaseSucceeded  StepPhase = "Succeeded"
	StepPhaseFailed     StepPhase = "Failed"
)

type StepStatus struct {
	StepName string
	Phase    StepPhase
	Reason   string
}

func NewStepStatus() *StepStatus {
	stepStatus := new(StepStatus)
	stepStatus.Phase = StepPhaseNotStarted
	return stepStatus
}

type LaunchPhase string

const (
	LaunchPhaseNotStarted LaunchPhase = "NotStarted"
	LaunchPhaseBuilding   LaunchPhase = "Processing"
	LaunchPhaseSucceeded  LaunchPhase = "Succeeded"
	LaunchPhaseFailed     LaunchPhase = "Failed"
)

type LaunchStatus struct {
	StepStatuses []StepStatus
	Phase        LaunchPhase
}

func NewLaunchStatus() *LaunchStatus {
	launchStatus := new(LaunchStatus)
	launchStatus.Phase = LaunchPhaseNotStarted
	return launchStatus
}
