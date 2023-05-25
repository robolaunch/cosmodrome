package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/robolaunch/cosmodrome/pkg/api"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type BuildLog struct {
	Stream string `json:"stream"`
}

func Build(ctx context.Context, dfName, dfPath, buildContext string, step api.Step, lc api.LaunchConfig) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	tar, err := archive.TarWithOptions(buildContext, &archive.TarOptions{})
	if err != nil {
		return err
	}

	buildArgs := make(map[string]*string)

	// make all keys upper case
	for k, v := range step.BuildArgs {
		buildArgs[strings.ToUpper(k)] = v
	}

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		tar,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dfPath + "/" + dfName,
			Remove:     true,
			Tags:       []string{step.Image.Name},
			BuildArgs:  buildArgs,
			NoCache:    !lc.Cache,
		},
	)

	if err != nil {
		return err
	}

	defer imageBuildResponse.Body.Close()

	if lc.Verbose {
		err = printBuildLogs(os.Stdout, imageBuildResponse.Body)
		if err != nil {
			return err
		}
	} else {
		var f *os.File
		if _, err := os.Stat("out_" + lc.Logfile); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir("out_"+lc.Logfile, os.ModePerm)
			if err != nil {
				return err
			}
			f, err = os.Create("out_" + lc.Logfile + "/out_" + step.Name + "_build_" + lc.Logfile)
			if err != nil {
				return err
			}
		} else {
			f, err = os.Create("out_" + lc.Logfile + "/out_" + step.Name + "_build_" + lc.Logfile)
			if err != nil {
				return err
			}
		}

		err = printBuildLogs(f, imageBuildResponse.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func BuildMultiplatform(ctx context.Context, dfName, dfPath, buildContext, baseImage string, step api.Step, lc api.LaunchConfig) error {

	// docker buildx rm multiarch_builder || true

	cmdDriverRemover := exec.CommandContext(
		context.Background(),
		"docker",
		"buildx",
		"rm",
		"multiarch_builder",
	)

	// docker buildx create --name multiarch_builder --use

	cmdDriverCreator := exec.CommandContext(
		context.Background(),
		"docker",
		"buildx",
		"create",
		"--name",
		"multiarch_builder",
		"--use",
	)

	platformStr := ""
	for _, platform := range step.Platforms {
		platformStr += platform + ","
	}
	platformStr = platformStr[0 : len(platformStr)-1]

	cmdBuildElements := []string{
		"buildx",
		"build",
		step.Context,
		"--file",
		step.Path + "/" + step.Dockerfile,
		"--platform",
		platformStr,
		"-t",
		step.Image.Name,
	}

	for k, v := range step.BuildArgs {
		cmdBuildElements = append(cmdBuildElements, "--build-arg")
		cmdBuildElements = append(cmdBuildElements, strings.ToUpper(k)+"="+*v)
	}

	if step.Push {
		cmdBuildElements = append(cmdBuildElements, "--push")
	}

	cmdBuilder := exec.CommandContext(
		context.Background(),
		"docker",
		cmdBuildElements...,
	)

	if lc.Verbose {
		cmdBuilder.Stdout = os.Stdout
		cmdBuilder.Stderr = os.Stderr
		cmdDriverRemover.Stdout = os.Stdout
		cmdDriverRemover.Stderr = os.Stderr
		cmdDriverCreator.Stdout = os.Stdout
		cmdDriverCreator.Stderr = os.Stderr
	} else {
		var f *os.File
		if _, err := os.Stat("out_" + lc.Logfile); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir("out_"+lc.Logfile, os.ModePerm)
			if err != nil {
				return err
			}
			f, err = os.Create("out_" + lc.Logfile + "/out_" + step.Name + "_build_" + lc.Logfile)
			if err != nil {
				return err
			}
		} else {
			f, err = os.Create("out_" + lc.Logfile + "/out_" + step.Name + "_build_" + lc.Logfile)
			if err != nil {
				return err
			}
		}

		cmdBuilder.Stdout = f
		cmdBuilder.Stderr = f
		cmdDriverRemover.Stdout = f
		cmdDriverRemover.Stderr = f
		cmdDriverCreator.Stdout = f
		cmdDriverCreator.Stderr = f
	}

	_ = cmdDriverRemover.Run()

	if err := cmdDriverCreator.Run(); err != nil {
		return err
	}
	if err := cmdBuilder.Run(); err != nil {
		return err
	}

	return nil
}

func printBuildLogs(dst io.Writer, rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()

		buildLogLine := &BuildLog{}
		json.Unmarshal([]byte(lastLine), buildLogLine)
		if buildLogLine.Stream != "" {
			fmt.Fprint(dst, buildLogLine.Stream)
		}

		errLine := &ErrorLine{}
		json.Unmarshal([]byte(lastLine), errLine)
		if errLine.Error != "" {
			fmt.Fprint(dst, errLine.Error)
			return errors.New(errLine.Error)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
