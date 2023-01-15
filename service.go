package main

import (
	"DockerProtector/lib/docker"
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/urfave/cli/v2"
)

func Service(c *cli.Context) error {
	docker.Events(func(eventsMsg events.Message){
		if eventsMsg.Type != "container" { return }
		if eventsMsg.Status != "start" { return }
		name := eventsMsg.Actor.Attributes["name"]
		fmt.Println(name)
	})
	return nil
}
