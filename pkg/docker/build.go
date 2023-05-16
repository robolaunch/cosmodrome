package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/robolaunch/cosmodrome/pkg/api"
)

func Build(ctx context.Context, dfName, dfPath, baseImage string, step api.Step, lc api.LaunchConfig) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	tar, err := archive.TarWithOptions(dfPath, &archive.TarOptions{})
	if err != nil {
		return err
	}

	buildArgs := make(map[string]*string)

	// make all keys upper case
	for k, v := range step.BuildArgs {
		buildArgs[strings.ToUpper(k)] = v
	}

	if baseImage != "" {
		buildArgs["BASE_IMAGE"] = &baseImage
	}

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		tar,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dfName,
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
		_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
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

		_, err = io.Copy(f, imageBuildResponse.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func BuildMultiplatform(ctx context.Context, dfName, dfPath, baseImage string, step api.Step, lc api.LaunchConfig) error {

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

	cmdBuildElements := []string{
		"buildx",
		"build",
		step.Path,
		"--file",
		step.Path + "/" + step.Dockerfile,
		"--platform",
		"linux/amd64,linux/arm64",
		"-t",
		step.Image.Name,
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
