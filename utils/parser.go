package utils

import (
	"fmt"
	"strings"
	"time"
)

type File struct {
	File    string `json:"file"`
	Content string `json:"content"`
}

func ExtractJSON(s string) string {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "```") {
		lines := strings.SplitN(s, "\n", 2)
		if len(lines) > 1 {
			s = lines[1]
		}
		if idx := strings.Index(s, "```"); idx != -1 {
			s = s[:idx]
		}
	}

	s = strings.TrimSpace(s)

	// Fix invalid JSON: AI may produce \$ (invalid escape) for Dart's $variable.
	// Replace with \u0024 (valid JSON for $). Preserves \\$ as \\u0024 -> \$
	s = strings.ReplaceAll(s, `\$`, `\u0024`)

	return s
}

func ExtractAppName(text string) string {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "app:") {
			name := strings.TrimSpace(strings.Split(line, ":")[1])
			name = strings.ReplaceAll(name, " ", "_")
			return name
		}
	}

	return fmt.Sprintf("ai_app_%d", time.Now().Unix())
}
