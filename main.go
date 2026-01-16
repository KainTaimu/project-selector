package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"bookmark/cli"
)

type flags struct {
	EditMode   bool
	AppendMode bool
}

func main() {
	flags := parseFlags()

	if flags.AppendMode {
		if err := appendMode(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}

	if flags.EditMode {
		if err := editMode(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}

	if flag.NArg() != 0 {
		if err := cli.RunQuickJumper(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err := cli.RunSelector(); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

// Appends the current working directory to the bookmarks file
func appendMode() (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %w", err)
	}
	pwd = cli.ShortenTildeExpansion(pwd)

	bookmarksFilePath := cli.GetBookmarksFilePath()

	file, err := os.OpenFile(bookmarksFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open bookmarks file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var fileLen int64
	if f, err := file.Stat(); err == nil {
		fileLen = f.Size()
	} else {
		return fmt.Errorf("failed to open config stat: %w", err)
	}

	// Add a newline to end of config file
	s := pwd
	buf := make([]byte, 1)
	n, err := file.ReadAt(buf, fileLen-1)
	if n != 0 {
		if err != nil {
			return fmt.Errorf("failed to read config file tail: %w", err)
		}
		if buf[0] != '\n' {
			s = "\n" + pwd
		}
	}

	_, err = file.Write([]byte(s))
	if err != nil {
		return fmt.Errorf("failed to append to config: %w", err)
	}

	return nil
}

// Runs an editor program like vim provided by environment variable $VISUAL or $EDITOR, before continuing to normal behaviour after the editor process is stopped.
// $EDITOR is used first if set, otherwise, $VISUAL is used.
//
// Returns an error if neither $EDITOR or $VISUAL is set.
//
// BUG(Erwin): Editors that run a detached process from the terminal like `code` from VSCode will cause bookmark to continue to
// normal behaviour without waiting for the user to save their changes to the config file.
func editMode() (err error) {
	var cmd exec.Cmd
	bookmarksFilePath := cli.GetBookmarksFilePath()

	if editor, exists := os.LookupEnv("EDITOR"); exists {
		cmd = *exec.Command(editor, bookmarksFilePath)
	} else if editor, exists := os.LookupEnv("VISUAL"); exists {
		cmd = *exec.Command(editor, bookmarksFilePath)
	} else {
		return fmt.Errorf("no editor available. Set $EDITOR or $VISUAL to your editor")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run editor: %w", err)
	}
	return nil
}

func parseFlags() flags {
	editMode := flag.Bool("e", false, "Launch an editor set by $EDITOR or $VISUAL")
	appendMode := flag.Bool("a", false, "Append current directory to bookmarks file")
	flag.Parse()

	flags := flags{
		*editMode,
		*appendMode,
	}

	return flags
}
