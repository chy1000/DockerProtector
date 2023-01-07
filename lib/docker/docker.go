package docker

import (
	"DockerProtector/lib/logger"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

var (
	ctx context.Context
	cli *client.Client
)

type Container struct {
	Name string // 容器名称
	Data types.ContainerJSON // 容器信息
}

func init(){
	var err error
	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal("执行 docker 客户端失败", err)
	}
}

func (c *Container) Inspect() error{
	container, err := cli.ContainerInspect(ctx, c.Name)
	if err != nil { return errors.New("获取容器信息失败") }
	c.Data = container
	return nil
}

func (c *Container) GetIp() (string, error) {
	if c.Data.NetworkSettings == nil {
		err := c.Inspect()
		if err != nil { return "", err }
	}
	networks := c.Data.NetworkSettings.Networks
	return networks["bridge"].IPAddress, nil
}

func Events(eventFunc func(events.Message)) {
	msgs, errs := cli.Events(ctx, types.EventsOptions{})
	for {
		select {
			case err := <-errs: fmt.Println(err)
			case msg := <-msgs:
				eventFunc(msg)
		}
	}
}