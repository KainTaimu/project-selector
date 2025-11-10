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

func IsFile(file string) (isFile bool) {
	if len(file) <= 0 {
		return false
	}

	if file[0] == '~' {
		file = os.Getenv("HOME") + file[1:]
	}

	if stat, err := os.Stat(file); err == nil {
		return !stat.IsDir() // Return true if file is not dir
	} else {
		return false
	}
}

func IsDir(file string) (isDir bool) {
	if len(file) <= 0 {
		return false
	}

	if file[0] == '~' {
		file = os.Getenv("HOME") + file[1:]
	}

	if stat, err := os.Stat(file); err == nil {
		return stat.IsDir()
	} else {
		return false
	}
}
