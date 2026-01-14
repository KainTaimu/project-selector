package cli

import (
	"os"
	"path/filepath"
	"strings"
)

// GetProjectsConfig returns the path to the main config at "~/.config/project_selector/projects.conf"
func GetProjectsConfig() string {
	return filepath.Join(os.Getenv(ConfigHomeEnv), AppConfigDir, ProjectEntriesFile)
}

// TildeExpansion attempts to expand the home "~/" string in the front of s into the path set by $HOME.
// Returns s as-is if "~/" is not the first two characters.
func TildeExpansion(s string) string {
	if strings.HasPrefix(s, "~/") {
		home := os.Getenv("HOME")
		s = filepath.Join(home, s[1:])
	}
	return s
}

// ShortenTildeExpansion attempts to shorten the path s from an absolute path to a relative path from $HOME
func ShortenTildeExpansion(s string) string {
	home := os.Getenv("HOME")
	if strings.HasPrefix(s, os.Getenv("HOME")) {
		s = filepath.Join("~", s[len(home):])
	}
	return s
}

func IsFile(file string) (isFile bool) {
	if len(file) <= 0 {
		return false
	}

	file = TildeExpansion(file)

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

	file = TildeExpansion(file)

	if stat, err := os.Stat(file); err == nil {
		return stat.IsDir()
	} else {
		return false
	}
}

// IsEmptyString returns true if s only consists of whitespace.
func IsEmptyString(s string) bool {
	for _, c := range s {
		if c != ' ' {
			return false
		}
	}
	return true
}
