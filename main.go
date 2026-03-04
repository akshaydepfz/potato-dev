package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"potato-dev/api"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	http.HandleFunc("/generate-app", api.GenerateApp)

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}
