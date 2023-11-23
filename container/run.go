package container

import (
	"context"
	"io"
	"os"

	"github.com/go-zoox/core-utils/fmt"

	"github.com/docker/docker/api/types"
	co "github.com/docker/docker/api/types/container"
	dContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// RunOptions is the configuration for running a container
type RunOptions struct {
	Name string
	//
	Container *co.Config
	Host      *co.HostConfig
	Network   *network.NetworkingConfig
	Platform  *specs.Platform
	//
	Detached bool
	//
	Stdin  io.Reader
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

// Run runs a container
func (c *container) Run(ctx context.Context, opts ...func(opt *RunOptions)) error {
	optsX := &RunOptions{
		Container: &co.Config{
			Tty:          true,
			OpenStdin:    true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			StdinOnce:    true,
		},
		Host:     &co.HostConfig{},
		Network:  &network.NetworkingConfig{},
		Platform: &specs.Platform{},
		//
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	for _, opt := range opts {
		opt(optsX)
	}

	if optsX.Detached {
		// optsX.Container.Tty = false
		// optsX.Container.OpenStdin = false
		optsX.Container.AttachStdin = false
		optsX.Container.AttachStdout = false
		optsX.Container.AttachStderr = false
		optsX.Container.StdinOnce = false
	}

	resp, err := c.Create(ctx, func(opt *CreateOptions) {
		opt.Name = optsX.Name
		opt.Container = optsX.Container
		opt.Host = optsX.Host
		opt.Network = optsX.Network
		opt.Platform = optsX.Platform
	})
	if err != nil {
		return err
	}
	containerID := resp.ID

	if !optsX.Detached {
		stream, err := c.client.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		})
		if err != nil {
			return err
		}
		go io.Copy(stream.Conn, optsX.Stdin)
		go io.Copy(optsX.Stdout, stream.Reader)

		// go func() {
		// 	for {
		// 		buf := make([]byte, 1024)
		// 		if n, err := stream.Reader.Read(buf); err != nil {
		// 			fmt.Println(err)
		// 		} else {
		// 			buf = buf[:n]

		// 			fmt.Println(buf)

		// 			// if buf[n-2] == '\r' && buf[n-1] == '\n' {
		// 			// 	continue
		// 			// }

		// 			// // handle ANI Ecape equence for cursor position => 27 91 54 110
		// 			// if buf[n-4] == 27 && buf[n-3] == 91 && buf[n-2] == 54 && buf[n-1] == 110 {
		// 			// 	buf = buf[:n-4]
		// 			// }

		// 			optsX.Stdout.Write(buf)
		// 		}
		// 	}
		// }()
	}

	if err := c.Start(ctx, resp.ID); err != nil {
		return err
	}

	if !optsX.Detached {
		statusCh, errCh := c.client.ContainerWait(ctx, resp.ID, dContainer.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			return err
		case status := <-statusCh:
			if status.StatusCode != 0 {
				return fmt.Errorf("container exited with non-zero status: %d", status.StatusCode)
			}
		}
	}

	return nil
}
