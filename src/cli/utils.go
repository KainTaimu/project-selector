package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetProjectsConfig() string {
	return os.Getenv(ConfigHomeEnv) + "/" + AppConfigDir + ProjectEntriesFile
}

func TildeExpansion(s string) (string, error) {
	if s[0] == '~' {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			return "", fmt.Errorf("$HOME is not set")
		}
		s = filepath.Join(home, s[1:])
		if !IsDir(s) {
			return "", fmt.Errorf("malformed working dir \"%s\". is $HOME set correctly?", s)
		}
	}
	return s, nil
}

func ShortenTildeExpansion(entry string) string {
	if home := os.Getenv("HOME"); strings.HasPrefix(entry, os.Getenv("HOME")) {
		entry = filepath.Join("~", entry[len(home):])
	}
	return entry
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
