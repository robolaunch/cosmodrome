package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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
		Password:      os.Getenv("REGISTRY_PAT"),
		ServerAddress: lc.Registry, // https://index.docker.io/v1/
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

	if lc.Verbose {
		_, err = io.Copy(os.Stdout, imagePushResponse)
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
			f, err = os.Create("out_" + lc.Logfile + "/out_" + step.Name + "_push_" + lc.Logfile)
			if err != nil {
				return err
			}
		} else {
			f, err = os.Create("out_" + lc.Logfile + "/out_" + step.Name + "_push_" + lc.Logfile)
			if err != nil {
				return err
			}
		}

		_, err = io.Copy(f, imagePushResponse)
		if err != nil {
			return err
		}

	}

	return nil
}
