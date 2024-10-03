package main

import (
	"dockermanager/dockerutils"
	"encoding/json"
	"log"
	"net/http"
)

func handleListContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := dockerutils.ListContainers()
	if err != nil {
		http.Error(w, "Error fetching containers", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(containers)
}

func handleStartContainer(w http.ResponseWriter, r *http.Request) {
	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		http.Error(w, "Missing container ID", http.StatusBadRequest)
		return
	}

	err := dockerutils.StartContainer(containerID)
	if err != nil {
		http.Error(w, "Error starting container", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Container started successfully"))
}

func handleStopContainer(w http.ResponseWriter, r *http.Request) {
	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		http.Error(w, "Missing container ID", http.StatusBadRequest)
		return
	}

	err := dockerutils.StopContainer(containerID)
	if err != nil {
		http.Error(w, "Error stopping container", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Container stopped successfully"))
}

func main() {
	http.HandleFunc("/api/containers", handleListContainers)
	http.HandleFunc("/api/start", handleStartContainer)
	http.HandleFunc("/api/stop", handleStopContainer)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
