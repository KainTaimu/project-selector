package cli

import (
	"fmt"
	"os"
	"strings"
)

const (
	HomeConfig         = "XDG_CONFIG_HOME"
	AppConfig          = "project_selector"
	ProjectEntriesFile = "projects.conf"
)

type Config struct {
	Projects []string
}

func ReadConfig() (entries []string, err error) {
	configDir := os.Getenv(HomeConfig) + "/"
	appConfig := configDir + AppConfig + "/"
	projectsPath := appConfig + ProjectEntriesFile

	err = os.MkdirAll(appConfig, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to create app config directory '%s': %w", appConfig, err)
	}

	if !IsFile(projectsPath) {
		_, err = os.Create(projectsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create projects config '%s': %w", projectsPath, err)
		}
	}

	file, err := os.ReadFile(projectsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects file '%s': %w", projectsPath, err)
	}

	contents := string(file)
	contents = strings.TrimSpace(contents)
	for entry := range strings.SplitSeq(contents, "\n") {
		entries = append(entries, entry)
	}

	return entries, nil
}

func IsFile(file string) (isFile bool) {
	if stat, err := os.Stat(file); err == nil {
		return !stat.IsDir() // Return true if file is not dir
	} else {
		return false
	}
}
