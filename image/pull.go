package image

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/cli/cli/streams"
	dimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/pkg/jsonmessage"
)

// PullConfig is the options for pulling an image
type PullConfig struct {
	Auth struct {
		Username string
		Password string
	}
	Platform string
	//
	Stdout io.Writer
}

// Pull pulls an image
func (i *image) Pull(ctx context.Context, name string, opts ...func(opt *PullConfig)) error {
	cfg := &PullConfig{
		Stdout: os.Stdout,
	}
	for _, o := range opts {
		o(cfg)
	}

	auth := ""
	if cfg.Auth.Username != "" && cfg.Auth.Password != "" {
		authConfig := registry.AuthConfig{
			Username: cfg.Auth.Username,
			Password: cfg.Auth.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		auth = base64.URLEncoding.EncodeToString(encodedJSON)
	}

	reader, err := i.client.ImagePull(ctx, name, dimage.PullOptions{
		RegistryAuth: auth,
		Platform:     cfg.Platform,
	})
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := jsonmessage.DisplayJSONMessagesToStream(reader, streams.NewOut(cfg.Stdout), nil); err != nil {
		return err
	}

	return nil
}
