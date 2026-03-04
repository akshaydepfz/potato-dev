package builder

import (
	"os"
	"os/exec"
	"path/filepath"

	"potato-dev/github"
	"potato-dev/utils"
)

func BuildProject(projectName string, files []utils.File, onStatus func(string)) (string, error) {
	status := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	projectDir := filepath.Join("workspace", projectName)

	status("Cleaning workspace...")
	os.RemoveAll(projectDir)

	status("Creating Flutter project...")
	cmd := exec.Command("flutter", "create", projectDir)
	cmd.Run()

	for _, file := range files {
		status("Creating file: " + file.File)

		filePath := filepath.Join(projectDir, file.File)

		dir := filepath.Dir(filePath)

		os.MkdirAll(dir, os.ModePerm)

		os.WriteFile(filePath, []byte(file.Content), 0644)
	}

	status("Creating GitHub repo...")
	repoURL, err := github.CreateAndPush(projectName, projectDir, status)
	if err != nil {
		return "", err
	}

	status("Done")
	return repoURL, nil
}
