package main

import (
	"fmt"
	"log"
	"os"

	docker "github.com/joshjms/inbox/client"
	"github.com/joshjms/inbox/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	// inbox run --dir=/path/to/dir

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// 	return
	// }

	app := cli.NewApp()
	app.Name = "inbox"
	app.Usage = "Running executables in a container"
	app.Commands = []*cli.Command{
		{
			Name:        "run",
			Description: "Run an executable in a container",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "pull",
					Usage: "Pull the image before running",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Println("Specify a file to run (e.g. app.exe)")

					return nil
				}
				p := c.Args().Get(0)
				withPull := c.Bool("pull")
				err := docker.Run(p, withPull)
				return utils.HandleError(err)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
