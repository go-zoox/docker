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

// PushConfig is the configuration for pushing an image
type PushConfig struct {
	Auth struct {
		Username string
		Password string
		Server   string
	}

	Stdout io.Writer
}

// Push pushes an image
func (i *image) Push(ctx context.Context, name string, opts ...func(cfg *PushConfig)) error {
	cfg := &PushConfig{
		Stdout: os.Stdout,
	}
	for _, o := range opts {
		o(cfg)
	}

	auth := "no-auth" // @TODO fix: bad parameters and missing X-Registry-Auth: invalid X-Registry-Auth header: EOF
	if cfg.Auth.Username != "" && cfg.Auth.Password != "" {
		authConfig := registry.AuthConfig{
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

	reader, err := i.client.ImagePush(ctx, name, dimage.PushOptions{
		All:          true,
		RegistryAuth: auth,
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
