package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigHomeEnv       = "XDG_CONFIG_HOME"
	AppConfigDir        = "bookmark/"
	BookmarkEntriesFile = "bookmarks.conf"
)

type Entry struct {
	Path    string
	IsValid bool
}

// ReadConfig reads the config in "~/.config/bookmark/" and produces a list of Entries from it
func ReadConfig() (entries []Entry, err error) {
	configDir := os.Getenv(ConfigHomeEnv) + "/"
	appConfig := filepath.Join(configDir, AppConfigDir)
	bookmarksFilePath := filepath.Join(appConfig, BookmarkEntriesFile)

	err = os.MkdirAll(appConfig, 0o777)
	if err != nil {
		return nil, fmt.Errorf("failed to create app config directory '%s': %w", appConfig, err)
	}

	if !IsFile(bookmarksFilePath) {
		_, err = os.Create(bookmarksFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create bookmarks config '%s': %w", bookmarksFilePath, err)
		}
	}

	file, err := os.ReadFile(bookmarksFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read bookmarks file '%s': %w", bookmarksFilePath, err)
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
