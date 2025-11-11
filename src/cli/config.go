package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigHomeEnv      = "XDG_CONFIG_HOME"
	AppConfigDir       = "project_selector/"
	ProjectEntriesFile = "projects.conf"
)

type Config struct {
	Projects []string
}

func ReadConfig() (entries []string, err error) {
	configDir := os.Getenv(ConfigHomeEnv) + "/"
	appConfig := filepath.Join(configDir, AppConfigDir)
	projectsPath := filepath.Join(appConfig, ProjectEntriesFile)

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
		entry = strings.TrimSpace(entry)
		if strings.HasPrefix(entry, "#") {
			continue
		}
		if IsEmptyString(entry) {
			continue
		}

		entries = append(entries, entry)
	}

	if missing := verifyProjectsExists(entries); missing != "" {
		return nil, fmt.Errorf("invalid entry found: directory '%s' does not exist", missing)
	}

	return entries, nil
}

func verifyProjectsExists(paths []string) string {
	for _, entry := range paths {
		if !IsDir(entry) {
			return entry
		}
	}
	return ""
}
