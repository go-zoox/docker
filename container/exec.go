package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
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

// Exec executes a command inside a container.
func (c *container) Exec(ctx context.Context, id string, opts ...func(opt *ExecOptions)) (io.ReadWriteCloser, error) {
	opt := &ExecOptions{
		// Stdin:  os.Stdin,
		// Stdout: os.Stdout,
		// Stderr: os.Stderr,
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
		return nil, err
	}

	stream, err := c.client.ContainerExecAttach(ctx, response.ID, types.ExecStartCheck{
		Detach: opt.Detach,
		Tty:    opt.Tty,
		// ConsoleSize: &[2]uint{300, 600},
	})
	if err != nil {
		return nil, err
	}

	// defer stream.Close()

	// if opt.Stdin != nil {
	// 	go io.Copy(stream.Conn, os.Stdin)
	// }

	// if _, err := io.Copy(opt.Stdout, stream.Reader); err != nil {
	// 	return err
	// }

	// return nil

	return &ExecStream{stream}, nil
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
