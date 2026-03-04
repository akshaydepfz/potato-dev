package ai

import (
	"context"
	"os"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

func GenerateUI(promptInput string) (string, error) {

	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	prompt := `
You are a mobile UI designer.

Generate a beautiful mobile app UI layout using HTML and Tailwind CSS.

User requirement:
` + promptInput + `

Rules:

1. Design for mobile width 375px
2. Use Tailwind CSS classes
3. Include modern spacing, shadows and rounded cards
4. Include header section
5. Include search bar
6. Include card components
7. Include bottom navigation
8. UI must look like a real mobile app
9. Use colors and icons where appropriate

Important:

Return ONLY valid HTML.

Do not include explanation.
Do not include markdown.
Do not include markdown code blocks (triple backticks).

Example structure:

<div class="w-[375px] mx-auto bg-gray-100 min-h-screen">
   header
   content
   bottom navigation
</div>
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
		return "", err
	}

	html := resp.OutputText()

	return html, nil
}