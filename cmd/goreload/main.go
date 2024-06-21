// Command interface for running the hot reload server
package main

import (
	"os"

	"github.com/idreaminteractive/goreload/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	var url string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url",
				Value:       "http://localhost:8082",
				Usage:       "Url where the reload server will be run",
				Destination: &url,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Runs hot reload server on a specific port",
				Action: func(cCtx *cli.Context) error {

					return commands.Serve(cCtx, url)
				},
			},
			{
				Name:  "reload",
				Usage: "Triggers hot reload",
				Action: func(cCtx *cli.Context) error {
					return commands.SignalReload(url)
				},
			},
		},
	}

	app.Run(os.Args)

}
