package github

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func CreateRepo(repoName string) error {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		return fmt.Errorf("GitHub token missing: ensure Koyeb secret GITHUB_TOKEN is set")
	}
	body := fmt.Sprintf(`{"name":"%s","private":false}`, repoName)

	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/user/repos",
		strings.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func Push(projectDir string, repoName string) error {
	username := strings.TrimSpace(os.Getenv("GITHUB_USERNAME"))
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if username == "" || token == "" {
		return fmt.Errorf("GitHub credentials missing: ensure Koyeb secrets GITHUB_USERNAME and GITHUB_TOKEN are set")
	}

	repoURL := fmt.Sprintf(
		"https://%s:%s@github.com/%s/%s.git",
		username,
		token,
		username,
		repoName,
	)

	cmds := [][]string{
		{"git", "init"},
		{"git", "add", "."},
		{"git", "commit", "-m", "AI generated Flutter project"},
		{"git", "branch", "-M", "main"},
		{"git", "remote", "add", "origin", repoURL},
		{"git", "push", "-u", "origin", "main", "--force"},
	}

	for _, c := range cmds {
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func RepoURL(repoName string) string {
	username := strings.TrimSpace(os.Getenv("GITHUB_USERNAME"))
	if username == "" {
		return ""
	}
	return fmt.Sprintf("https://github.com/%s/%s", username, repoName)
}

func CreateAndPush(projectName string, projectDir string, onStatus func(string)) (string, error) {
	status := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	if err := CreateRepo(projectName); err != nil {
		return "", err
	}

	status("Pushing to GitHub...")
	if err := Push(projectDir, projectName); err != nil {
		return "", err
	}
	return RepoURL(projectName), nil
}
