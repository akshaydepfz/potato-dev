package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"

	"potato-dev/utils"
)

func Generate(requirement string, onStatus func(string)) ([]utils.File, error) {
	status := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	status("Generating app with AI...")
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key missing: ensure Koyeb secret OPENAI_API_KEY is set")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	prompt := `
You are a senior Flutter architect.

Generate a COMPLETE production-ready Flutter mobile application.

Requirement:
` + requirement + `

Use Flutter best practices and clean architecture.

The generated project MUST follow this folder structure:

lib/
  main.dart
  app.dart

  core/
    theme/app_theme.dart
    constants/app_colors.dart
    constants/app_strings.dart
    utils/logger.dart

  models/
    user_model.dart
    product_model.dart
    order_model.dart

  services/
    firebase_service.dart
    auth_service.dart
    firestore_service.dart
    storage_service.dart

  providers/
    auth_provider.dart
    cart_provider.dart
    product_provider.dart

  routes/
    app_routes.dart

  screens/
    splash_screen.dart
    login_screen.dart
    register_screen.dart
    home_screen.dart
    cart_screen.dart
    profile_screen.dart

  widgets/
    custom_button.dart
    custom_textfield.dart
    product_card.dart
    loading_widget.dart

If the requirement mentions backend or authentication:

Use Firebase and include:

firebase_core
firebase_auth
cloud_firestore
firebase_storage

Initialize Firebase in main.dart.

State management must use Provider.

Each screen must include a functional UI layout.

Use modern Material 3 UI design.

Use dummy data where backend logic is incomplete.

Important Rules:

1. Return ONLY valid JSON
2. No explanations
3. Each file must contain complete Dart code
4. Imports must be correct
5. The project must compile without syntax errors

JSON response format must be EXACTLY:

[
  {
    "file": "lib/main.dart",
    "content": "dart code"
  },
  {
    "file": "lib/screens/home_screen.dart",
    "content": "dart code"
  }
]

Return ONLY JSON.
`

	resp, err := client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model: responses.ChatModelGPT4o,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: param.NewOpt(prompt),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	output := resp.OutputText()

	var files []utils.File

	err = json.Unmarshal([]byte(utils.ExtractJSON(output)), &files)

	if err != nil {
		return nil, err
	}

	return files, nil
}
