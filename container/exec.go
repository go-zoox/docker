package container

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

// ExecOptions is the options for Exec.
type ExecOptions struct {
	Detach bool
	Tty    bool
	Cmd    []string
	//
	Stdin  io.Reader
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

// Exec executes a command inside a container.
func (c *container) Exec(ctx context.Context, id string, opts ...func(opt *ExecOptions)) error {
	opt := &ExecOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	for _, o := range opts {
		o(opt)
	}

	response, err := c.client.ContainerExecCreate(ctx, id, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          opt.Tty,
		Cmd:          opt.Cmd,
		// ConsoleSize:  &[2]uint{300, 600},
	})
	if err != nil {
		return err
	}

	stream, err := c.client.ContainerExecAttach(ctx, response.ID, types.ExecStartCheck{
		Detach: opt.Detach,
		Tty:    opt.Tty,
		// ConsoleSize: &[2]uint{300, 600},
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	if opt.Stdin != nil {
		go io.Copy(stream.Conn, os.Stdin)
	}

	io.Copy(os.Stdout, stream.Reader)

	return nil
}
