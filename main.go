package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"potato-dev/api"
)

func main() {
	_ = godotenv.Load() // optional for local .env; Koyeb uses secrets

	// Validate required Koyeb secrets at startup
	required := map[string]string{
		"OPENAI_API_KEY":  strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
		"GITHUB_TOKEN":    strings.TrimSpace(os.Getenv("GITHUB_TOKEN")),
		"GITHUB_USERNAME": strings.TrimSpace(os.Getenv("GITHUB_USERNAME")),
		"TELEGRAM_TOKEN":  strings.TrimSpace(os.Getenv("TELEGRAM_TOKEN")),
	}
	for name, val := range required {
		if val == "" {
			log.Fatalf("Missing Koyeb secret %s: ensure it is set in Koyeb environment", name)
		}
	}

	http.HandleFunc("/generate-app", api.GenerateApp)

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}
