package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func Run(dir string) error {
	fmt.Println("Running at", dir)

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "busybox:stable-uclibc", image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	return nil
}
