package process

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var LaunchLog *color.Color = color.New(color.FgYellow).Add(color.Italic).Add(color.Bold)
var StepLog *color.Color = color.New(color.FgBlue)
var SuccessLog *color.Color = color.New(color.FgGreen).Add(color.Bold)
var logSpinner *spinner.Spinner = spinner.New(spinner.CharSets[4], 100*time.Millisecond)

func GetSpinner(logStyle *color.Color, msg string) {
	logSpinner.Suffix = logStyle.Sprintln(msg)
}
