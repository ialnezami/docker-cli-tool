package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

func startContainer(containerID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container %s started\n", containerID)
}

var startCmd = &cobra.Command{
	Use:   "start [container_id]",
	Short: "Start a Docker container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startContainer(args[0])
	},
}

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

func monitorContainerStats(containerID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	stats, err := cli.ContainerStats(context.Background(), containerID, true)
	if err != nil {
		panic(err)
	}

	// Handle stats streaming and display CPU, memory, and network usage
	//  Use decoder to parse stats.Body continuously
	//  Print CPU, memory, and network usage
	//  Use time.Sleep to control the frequency of stats updates
	//  Use stats.Body.Close() to close the stream
	//  Use stats.Body to read the stream
	//  Use stats.Body.Read() to read the stream

}

var statsCmd = &cobra.Command{
	Use:   "stats [container_id]",
	Short: "Monitor container resource usage",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		monitorContainerStats(args[0])
	},
}

func getContainerLogs(containerID string, follow bool) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow}
	out, err := cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		panic(err)
	}

	// Read the output and print logs to console
	io.Copy(os.Stdout, out)
}

var logsCmd = &cobra.Command{
	Use:   "logs [container_id]",
	Short: "Fetch logs from a Docker container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getContainerLogs(args[0], true)
	},
}

func createAndRunContainer(image string, cpuLimit float64, memoryLimit int64) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			CPUQuota: int64(cpuLimit * 100000),
			Memory:   memoryLimit * 1024 * 1024,
		},
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
	}, hostConfig, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// Start the container
	cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
}

var runCmd = &cobra.Command{
	Use:   "run [image]",
	Short: "Run a Docker container with resource limits",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cpuLimit, _ := cmd.Flags().GetFloat64("cpu")
		memoryLimit, _ := cmd.Flags().GetInt64("memory")
		createAndRunContainer(args[0], cpuLimit, memoryLimit)
	},
}

func init() {
	runCmd.Flags().Float64("cpu", 0, "Set CPU limit")
	runCmd.Flags().Int64("memory", 0, "Set memory limit in MB")
}

func getContainerHealth(containerID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containerJSON, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Health Status: %s\\n", containerJSON.State.Health.Status)
}

var healthCmd = &cobra.Command{
	Use:   "health [container_id]",
	Short: "Check the health status of a Docker container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getContainerHealth(args[0])
	},
}

func listImages() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Printf("ID: %s | RepoTags: %v\\n", image.ID[:10], image.RepoTags)
	}
}

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List Docker images",
	Run: func(cmd *cobra.Command, args []string) {
		listImages()
	},
}

func runDockerComposeUp() {
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

var composeUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run docker-compose up",
	Run: func(cmd *cobra.Command, args []string) {
		runDockerComposeUp()
	},
}

func execInContainer(containerID, cmd string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{cmd},
	}

	execID, err := cli.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerExecAttach(context.Background(), execID.ID, execConfig)
	if err != nil {
		panic(err)
	}

	defer resp.Close()
	io.Copy(os.Stdout, resp.Reader)
}

var execCmd = &cobra.Command{
	Use:   "exec [container_id] [cmd]",
	Short: "Run a command inside a Docker container",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		execInContainer(args[0], args[1])
	},
}

func listVolumes() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	volumes, err := cli.VolumeList(context.Background(), filters.Args{})
	if err != nil {
		panic(err)
	}

	for _, volume := range volumes.Volumes {
		fmt.Printf("Volume Name: %s\\n", volume.Name)
	}
}

var volumeCmd = &cobra.Command{
	Use:   "volume ls",
	Short: "List Docker volumes",
	Run: func(cmd *cobra.Command, args []string) {
		listVolumes()
	},
}
