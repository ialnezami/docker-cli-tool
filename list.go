package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func listContainers() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("ID: %s | Image: %s | State: %s\n", container.ID[:10], container.Image, container.State)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Docker containers",
	Run: func(cmd *cobra.Command, args []string) {
		listContainers()
	},
}
