package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func stopContainer(containerID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container %s stopped\n", containerID)
}

var stopCmd = &cobra.Command{
	Use:   "stop [container_id]",
	Short: "Stop a Docker container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stopContainer(args[0])
	},
}
