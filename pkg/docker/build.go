package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/robolaunch/cosmodrome/pkg/api"
)

func Build(ctx context.Context, dfName, dfPath, baseImage string, step api.Step) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFileReader, err := os.Open(dfPath)
	if err != nil {
		return err
	}
	readDockerFile, err := io.ReadAll(dockerFileReader)
	if err != nil {
		return err
	}

	tarHeader := &tar.Header{
		Name: dfName,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return err
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return err
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	buildArgs := make(map[string]*string)
	if baseImage != "" {
		buildArgs["BASE_IMAGE"] = &baseImage
	}

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dfName,
			Remove:     true,
			Tags:       []string{step.Image.Name},
			BuildArgs:  buildArgs,
		},
	)

	if err != nil {
		return err
	}

	defer imageBuildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}
