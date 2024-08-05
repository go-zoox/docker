package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/api/types/container"
)

// ExecOptions is the options for Exec.
type ExecOptions struct {
	Detach bool
	Tty    bool
	Cmd    []string
	// //
	// Stdin  io.Reader
	// Stdout io.WriteCloser
	// Stderr io.WriteCloser
}

type ExecTerm struct {
	io.ReadWriteCloser
	Resize func(width, height uint) error
}

// Exec executes a command inside a container.
func (c *container) Exec(ctx context.Context, id string, opts ...func(opt *ExecOptions)) (*ExecTerm, error) {
	opt := &ExecOptions{
		// Stdin:  os.Stdin,
		// Stdout: os.Stdout,
		// Stderr: os.Stderr,
	}
	for _, o := range opts {
		o(opt)
	}

	response, err := c.client.ContainerExecCreate(ctx, id, dc.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          opt.Tty,
		Cmd:          opt.Cmd,
		// ConsoleSize:  &[2]uint{300, 600},
	})
	if err != nil {
		return nil, err
	}

	stream, err := c.client.ContainerExecAttach(ctx, response.ID, dc.ExecStartOptions{
		Detach: opt.Detach,
		Tty:    opt.Tty,
		// ConsoleSize: &[2]uint{300, 600},
	})
	if err != nil {
		return nil, err
	}

	// if err := c.client.ContainerExecStart(ctx, response.ID, dc.ExecStartOptions{}); err != nil {
	// 	return nil, nil, err
	// }

	resize := func(width, height uint) error {
		if !opt.Tty {
			return nil
		}

		return c.client.ContainerExecResize(ctx, response.ID, dc.ResizeOptions{
			Width:  width,
			Height: height,
		})
	}

	// defer stream.Close()

	// if opt.Stdin != nil {
	// 	go io.Copy(stream.Conn, os.Stdin)
	// }

	// if _, err := io.Copy(opt.Stdout, stream.Reader); err != nil {
	// 	return err
	// }

	// return nil

	return &ExecTerm{
		ReadWriteCloser: &ExecStream{stream},
		Resize:          resize,
	}, nil
}

type ExecStream struct {
	types.HijackedResponse
}

func (s *ExecStream) Close() error {
	s.HijackedResponse.Close()
	return nil
}

func (s *ExecStream) Read(p []byte) (n int, err error) {
	return s.Reader.Read(p)
}

func (s *ExecStream) Write(p []byte) (n int, err error) {
	return s.Conn.Write(p)
}
