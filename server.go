package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ID    string `json:"id"`
	Names string `json:"names"`
	Image string `json:"image"`
	State string `json:"state"`
}

func main() {
	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	// Endpoint to fetch container list
	http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			http.Error(w, "Error fetching containers", http.StatusInternalServerError)
			return
		}

		containerInfos := []ContainerInfo{}
		for _, container := range containers {
			containerInfos = append(containerInfos, ContainerInfo{
				ID:    container.ID,
				Names: container.Names[0],
				Image: container.Image,
				State: container.State,
			})
		}

		json.NewEncoder(w).Encode(containerInfos)
	})

	// Endpoint to fetch container logs
	http.HandleFunc("/container-logs", func(w http.ResponseWriter, r *http.Request) {
		containerID := r.URL.Query().Get("id")
		if containerID == "" {
			http.Error(w, "Missing container ID", http.StatusBadRequest)
			return
		}

		options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
		logs, err := cli.ContainerLogs(context.Background(), containerID, options)
		if err != nil {
			http.Error(w, "Error fetching logs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.ReadFrom(logs)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
