package image

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// PushOption is the configuration for pushing an image
type PushOption struct {
	Name string
	Auth struct {
		Username string
		Password string
		Server   string
	}
}

// Push pushes an image
func Push(cfg *PushOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	auth := ""
	if cfg.Auth.Username != "" && cfg.Auth.Password != "" {
		authConfig := types.AuthConfig{
			Username:      cfg.Auth.Username,
			Password:      cfg.Auth.Password,
			ServerAddress: cfg.Auth.Server,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		auth = base64.URLEncoding.EncodeToString(encodedJSON)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	reader, err := cli.ImagePush(ctx, cfg.Name, types.ImagePushOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	return nil
}
