package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/robolaunch/cosmodrome/pkg/api"
)

func Build(ctx context.Context, dfName, dfPath, baseImage string, step api.Step) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("unable to init client")
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFileReader, err := os.Open(dfPath)
	if err != nil {
		return errors.New("unable to open Dockerfile")
	}
	readDockerFile, err := io.ReadAll(dockerFileReader)
	if err != nil {
		return errors.New("unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dfName,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return errors.New("unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return errors.New("unable to write tar body")
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dfName,
			Remove:     true,
			Tags:       []string{step.Image.Name},
			BuildArgs: map[string]*string{
				"BASE_IMAGE": &baseImage,
			},
		},
	)

	if err != nil {
		return errors.New("unable to build docker image")
	}

	defer imageBuildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return errors.New("unable to read image build response")
	}

	return nil
}
