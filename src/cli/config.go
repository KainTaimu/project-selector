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

func GetProjectsConfig() string {
	return os.Getenv(ConfigHomeEnv) + "/" + AppConfigDir + ProjectEntriesFile
}

func ShortenTildeExpansion(entry string) string {
	if home := os.Getenv("HOME"); strings.HasPrefix(entry, os.Getenv("HOME")) {
		entry = filepath.Join("~", entry[len(home):])
	}
	return entry
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

func IsEmptyString(s string) bool {
	for _, c := range s {
		if c != ' ' {
			return false
		}
	}
	return true
}
