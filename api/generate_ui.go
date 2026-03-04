package api

import (
	"encoding/json"
	"fmt"
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

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.URL.Query().Get("preview") == "1" {
		// Return full HTML preview page with phone frame
		htmlJSON, _ := json.Marshal(html)
		previewHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-200 flex justify-center p-10">
<div class="bg-black p-3 rounded-[40px] shadow-xl">
  <div id="phone" class="bg-white rounded-[30px] overflow-hidden w-[375px]"></div>
</div>
<script>
document.getElementById("phone").innerHTML = %s;
</script>
</body>
</html>`, string(htmlJSON))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(previewHTML))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"html": html})
}
