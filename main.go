package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name: "DockerProtector",
		Usage: "docker protector service",
		Commands: []*cli.Command{
			{
				Name:  "cmd",
				Usage: "command",
				Action: Cmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "list",
						Aliases:  []string{"l"},
						Usage:    "显示容器的可访问IP的规则",
					},
					&cli.StringFlag{
						Name:     "add",
						Aliases:  []string{"a"},
						Usage:    "添加容器的可访问IP规则",
					},
					&cli.StringFlag{
						Name:     "remove",
						Aliases:  []string{"r"},
						Usage:    "删除容器的可访问IP规则",
					},
				},
			},
		},
		Action: Service,
	}
}

func main() {
	_ = app.Run(os.Args)
}