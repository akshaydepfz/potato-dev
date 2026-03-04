package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"potato-dev/ai"
)

type GenerateUIRequest struct {
	Prompt string `json:"prompt"`
}

func GenerateUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		http.Error(w, "OpenAI API key missing: ensure Koyeb secret OPENAI_API_KEY is set", http.StatusInternalServerError)
		return
	}

	var req GenerateUIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		http.Error(w, "prompt is required", http.StatusBadRequest)
		return
	}

	html, err := ai.GenerateUI(req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{"html": html})
}
