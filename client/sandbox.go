package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/joshjms/inbox/utils"
)

type Sandbox struct {
	ID         string
	Executable string
	Client     *DockerClient
}

func (s *Sandbox) Run() error {
	err := s.Init()
	if err != nil {
		return err
	}

	err = s.RunContainer()
	if err != nil {
		return err
	}

	return nil
}

func NewSandbox(e string, cli *DockerClient) *Sandbox {
	return &Sandbox{
		ID:         uuid.New().String(),
		Executable: e,
		Client:     cli,
	}
}

func (s *Sandbox) Init() error {
	os.Mkdir(s.ID, 0755)
	os.Create(filepath.Join(s.ID, "app"))
	err := utils.Copy(s.Executable, filepath.Join(s.ID, "app"))
	if err != nil {
		log.Println(err)
		return err
	}

	os.Create(filepath.Join(s.ID, "stdin.txt"))
	os.Create(filepath.Join(s.ID, "stdout.txt"))
	os.Create(filepath.Join(s.ID, "stderr.txt"))

	return nil
}

func (s *Sandbox) RunContainer() error {
	ctx := context.Background()

	resp, err := s.Client.ContainerCreate(ctx,
		&container.Config{
			Image: "busybox:latest",
			Cmd:   []string{"chmod", "+x", "/app/app", "&&", "/app/app"},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: filepath.Join(os.Getenv("MOUNTS_DIR"), s.ID),
					Target: "/app",
				},
			},
		},
		nil,
		nil,
		"",
	)

	if err != nil {
		return utils.HandleError(err)
	}

	if err := s.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return utils.HandleError(err)
	}

	fmt.Println("waiting")

	statusCh, errCh := s.Client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return utils.HandleError(err)
		}
	case <-statusCh:
	}

	fmt.Println("done")

	out, err := s.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		fmt.Println(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}
