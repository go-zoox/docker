package container

import (
	"context"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Logs retrieves the logs of a container
func Logs(id string) ([]byte, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	response, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return nil, err
	}
	defer response.Close()

	// It is correct
	// stdcopy.StdCopy(os.Stdout, os.Stderr, response)

	// https://stackoverflow.com/questions/39540079/explanation-of-docker-attach-payload
	ignoreHeader := make([]byte, 8)
	response.Read(ignoreHeader)

	// but this is not correct with error chars
	// io.Copy(os.Stdout, response)

	bytes, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
