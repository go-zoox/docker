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

// PullOption is the options for pulling an image
type PullOption struct {
	Auth struct {
		Username string
		Password string
	}
	Platform string
	//
	Stdout io.Writer
}

// Pull pulls an image
func (i *image) Pull(ctx context.Context, name string, opts ...func(opt *PullOption)) error {
	opt := &PullOption{
		Stdout: os.Stdout,
	}
	for _, o := range opts {
		o(opt)
	}

	auth := ""
	if opt.Auth.Username != "" && opt.Auth.Password != "" {
		authConfig := registry.AuthConfig{
			Username: opt.Auth.Username,
			Password: opt.Auth.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		auth = base64.URLEncoding.EncodeToString(encodedJSON)
	}

	reader, err := i.client.ImagePull(ctx, name, dimage.PullOptions{
		RegistryAuth: auth,
		Platform:     opt.Platform,
	})
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := jsonmessage.DisplayJSONMessagesToStream(reader, streams.NewOut(opt.Stdout), nil); err != nil {
		return err
	}

	return nil
}
