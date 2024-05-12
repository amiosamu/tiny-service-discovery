package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	HelloServiceImageName = "hello"
	ContainerRunningState = "running"
	ContainerKillState    = "kill"
	ContainerStartState   = "start"
)

type DockerClient struct {
	*client.Client
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &DockerClient{cli}, nil
}

func (dc *DockerClient) GetContainerPort(ctx context.Context, id string) (uint16, error) {
	containers, err := dc.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("id", id)),
		All:     true, // Include stopped containers to find the port mapping
	})
	if err != nil {
		return 0, err
	}

	for _, container := range containers {
		if container.ID == id {
			for _, port := range container.Ports {
				if port.PublicPort != 0 { // Ensure a public port is mapped
					return port.PublicPort, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("container %s not found or does not have a public port exposed", id)
}
