package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	*client.Client
}

func Run(dir string, withPull bool) error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	if withPull {
		if err := client.Pull(); err != nil {
			return err
		}
	}

	sandbox := NewSandbox(dir, client)
	fmt.Println(sandbox.ID)

	sandbox.Run()

	return nil
}

func NewClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli}, nil
}

func (cli DockerClient) Pull() error {
	ctx := context.Background()

	reader, err := cli.ImagePull(ctx, "busybox:latest", image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	return nil
}
