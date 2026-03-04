package main

import (
	"log"
	"net/http"

	"potato-dev/api"
)

func main() {

	// Secrets are validated when endpoints are called (not at startup)
	// In Koyeb: add env vars that reference secrets, e.g. OPENAI_API_KEY={{ secret.OPENAI_API_KEY }}

	http.HandleFunc("/generate-app", api.GenerateApp)
	http.HandleFunc("/generate-ui", api.GenerateUI)

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}
