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

type Entry struct {
	Path    string
	IsValid bool
}

// ReadConfig reads the config in "~/.config/project_selector/" and produces a list of Entries from it
func ReadConfig() (entries []Entry, err error) {
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

		if IsDir(entry) {
			entries = append(entries, Entry{Path: entry, IsValid: true})
		} else {
			entries = append(entries, Entry{Path: entry, IsValid: false})
		}
	}

	return entries, nil
}
