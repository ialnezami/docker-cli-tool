package main

import (
	"dockermanager/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all containers",
	Run: func(cmd *cobra.Command, args []string) {
		containers, err := utils.ListContainers()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		for _, container := range containers {
			fmt.Printf("Container ID: %s, Name: %s, State: %s\n", container.ID, container.Names[0], container.State)
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start [containerID]",
	Short: "Start a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]
		err := utils.StartContainer(containerID)
		if err != nil {
			log.Fatalf("Error starting container: %v", err)
		}
		fmt.Println("Container started successfully")
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [containerID]",
	Short: "Stop a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]
		err := utils.StopContainer(containerID)
		if err != nil {
			log.Fatalf("Error stopping container: %v", err)
		}
		fmt.Println("Container stopped successfully")
	},
}

func main() {
	var rootCmd = &cobra.Command{Use: "dockcli"}
	rootCmd.AddCommand(listCmd, startCmd, stopCmd)
	rootCmd.Execute()
}
