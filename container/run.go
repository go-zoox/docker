package container

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-zoox/docker/image"

	"github.com/docker/docker/api/types"
	dContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-zoox/fs"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/go-zoox/crypto/md5"
)

// RunConfig is the configuration for running a container
type RunConfig struct {
	Name       string
	Image      string
	Commands   []string
	Volumes    map[string]string
	WorkingDir struct {
		Host      string
		Container string
	}
	//
	Env map[string]string
	//
	Auth struct {
		Username string
		Password string
	}
	//
	Stdout io.Writer
	Stderr io.Writer
	//
	// Limit
	User   string
	Memory int64 // Memory limit (in bytes)
	CPU    int64 // CPU cores count

	// Dockerfile Content, not a file path or name
	Dockerfile string

	//
	AutoRemove bool
}

// Run runs a container
func Run(cfg *RunConfig) error {
	HostDir := cfg.WorkingDir.Host
	ContainerDir := cfg.WorkingDir.Container
	var cmd []string = nil

	// create tmp dir
	if err := fs.Mkdirp(HostDir); err != nil {
		return err
	}
	// // create log file
	// logf, err := os.Create(logPath)
	// if err != nil {
	// 	return err
	// }

	if cfg.Commands != nil {
		scriptHostPath := fs.JoinPath(HostDir, "runner.sh")
		scriptContainerPath := fs.JoinPath(ContainerDir, "runner.sh")

		// generates script by commands
		f, err := os.Create(scriptHostPath)
		if err != nil {
			return err
		}
		f.WriteString("#!/bin/sh\n")
		for _, cmd := range cfg.Commands {
			f.WriteString(cmd + "\n")
		}
		f.Close()

		cmd = []string{"sh", scriptContainerPath}
	}

	// logPath := fs.JoinPath(runtimeDir, "runner.log")

	//
	stdout := cfg.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := cfg.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	auth := ""
	if cfg.Auth.Username != "" && cfg.Auth.Password != "" {
		authConfig := types.AuthConfig{
			Username: cfg.Auth.Username,
			Password: cfg.Auth.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		auth = base64.URLEncoding.EncodeToString(encodedJSON)
	}

	env := []string{}
	for k, v := range cfg.Env {
		env = append(env, k+"="+v)
	}

	if cfg.Dockerfile != "" {
		dockerfilePath := fs.JoinPath(HostDir, "Dockerfile")
		// generates script by commands
		f, err := os.Create(dockerfilePath)
		if err != nil {
			return err
		}
		f.WriteString(cfg.Dockerfile)

		hash := md5.Md5(cfg.Dockerfile)
		imageName := strings.ToLower("GO_ZOOX_DOCKER_BUILD_" + hash + ":latest")

		err = image.Build(&image.BuildConfig{
			Name:    imageName,
			Context: HostDir,
			// Tags:    []string{"latest"},
		})
		if err != nil {
			return err
		}

		cfg.Image = imageName
	}

	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: HostDir,
			Target: ContainerDir,
		},
	}
	for k, v := range cfg.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: k,
			Target: v,
		})
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&dContainer.Config{
			Image:      cfg.Image,
			Cmd:        cmd,
			Env:        env,
			WorkingDir: cfg.WorkingDir.Container,
			User:       cfg.User,
		},
		&dContainer.HostConfig{
			Mounts: mounts,
			Resources: dContainer.Resources{
				Memory:   cfg.Memory,
				NanoCPUs: cfg.CPU,
			},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		cfg.Name,
	)
	if err != nil {
		if client.IsErrNotFound(err) {
			reader, err := cli.ImagePull(ctx, cfg.Image, types.ImagePullOptions{
				RegistryAuth: auth,
			})
			if err != nil {
				return err
			}
			io.Copy(stdout, reader)

			//
			resp, _ = cli.ContainerCreate(
				ctx,
				&dContainer.Config{
					Image:      cfg.Image,
					Cmd:        cmd,
					Env:        env,
					WorkingDir: cfg.WorkingDir.Container,
				},
				&dContainer.HostConfig{
					Mounts: mounts,
				},
				&network.NetworkingConfig{},
				&v1.Platform{},
				cfg.Name,
			)
		} else {
			return err
		}
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	go func(ctx context.Context, cli *client.Client, containerID string) {
		out, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Follow:     true,
		})
		if err != nil {
			return
		}

		// StdCopy is a modified version of io.Copy.
		stdcopy.StdCopy(stdout, stderr, out)
	}(ctx, cli, resp.ID)

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, dContainer.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return err
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with code %d", status.StatusCode)
		}
	}

	if cfg.AutoRemove {
		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
	}

	return nil
}
