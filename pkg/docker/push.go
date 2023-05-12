package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/robolaunch/cosmodrome/pkg/api"
)

func Push(ctx context.Context, step api.Step, lc api.LaunchConfig) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	var authConfig = types.AuthConfig{
		Username:      lc.Organization,
		Password:      "<PERSONAL-ACCESS-TOKEN>",
		ServerAddress: "<REGISTRY-SERVER-ADDRESS>", // https://index.docker.io/v1/
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	imagePushResponse, err := cli.ImagePush(
		ctx,
		step.Image.Name,
		types.ImagePushOptions{
			RegistryAuth: authConfigEncoded,
		},
	)

	if err != nil {
		return err
	}

	defer imagePushResponse.Close()

	_, err = io.Copy(os.Stdout, imagePushResponse)
	if err != nil {
		return err
	}

	return nil
}
