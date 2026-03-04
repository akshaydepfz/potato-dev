package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"potato-dev/ai"
	"potato-dev/builder"
)

type GenerateRequest struct {
	AppName     string `json:"app_name"`
	Requirement string `json:"requirement"`
}

type GenerateResponse struct {
	RepoURL string `json:"repo_url"`
	Status  string `json:"status"`
}

type StatusEvent struct {
	Status  string `json:"status"`
	RepoURL string `json:"repo_url,omitempty"`
	Error   string `json:"error,omitempty"`
}

func GenerateApp(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Stream status updates via Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	statusChan := make(chan string, 32)
	resultChan := make(chan struct {
		repoURL string
		err     error
	}, 1)

	go func() {
		onStatus := func(s string) {
			statusChan <- s
		}

		files, err := ai.Generate(req.Requirement, onStatus)
		if err != nil {
			resultChan <- struct {
				repoURL string
				err     error
			}{"", err}
			return
		}

		repoURL, err := builder.BuildProject(req.AppName, files, onStatus)
		resultChan <- struct {
			repoURL string
			err     error
		}{repoURL, err}
	}()

	for {
		select {
		case status := <-statusChan:
			event := StatusEvent{Status: status}
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

		case result := <-resultChan:
			if result.err != nil {
				event := StatusEvent{Status: "error", Error: result.err.Error()}
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", data)
			} else {
				event := StatusEvent{Status: "success", RepoURL: result.repoURL}
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", data)
			}
			flusher.Flush()
			return
		}
	}
}
