package main

import (
	"fmt"
	"log"
	"os"

	docker "github.com/joshjms/inbox/client"
	"github.com/urfave/cli/v2"
)

func main() {
	// inbox run --dir=/path/to/dir

	app := cli.NewApp()
	app.Name = "inbox"
	app.Usage = "Running executables in a container"
	app.Commands = []*cli.Command{
		{
			Name: "run",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Println("Specify a file to run (e.g. app.exe)")

					return nil
				}
				p := c.Args().Get(0)
				err := docker.Run(p)
				return err
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
