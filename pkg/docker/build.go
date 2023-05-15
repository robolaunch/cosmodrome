package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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
