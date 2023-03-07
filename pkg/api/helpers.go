package api

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func execBashCmd(bashCommand string) {
	cmd := exec.Command("/bin/bash", "-c", bashCommand)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	go getStd(cmd, stderr)
	go getStd(cmd, stdout)

	cmd.Wait()
}

func getStd(cmd *exec.Cmd, std io.ReadCloser) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

}

// func getStderr() {

// }
