package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"project-selector/src/cli"
)

type flags struct {
	EditMode   bool
	AppendMode bool
}

func main() {
	flags := parseFlags()

	if flags.EditMode {
		if err := editMode(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}

	if flags.AppendMode {
		if err := appendMode(); err != nil {
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

func appendMode() (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %w", err)
	}
	if home := os.Getenv("HOME"); strings.HasPrefix(pwd, os.Getenv("HOME")) {
		pwd = filepath.Join("~", pwd[len(home):])
	}
	pwd = pwd + "/"

	projectsFilePath := cli.GetProjectsConfig()

	file, err := os.OpenFile(projectsFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open config: %w", err)
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

	var s string
	buf := make([]byte, 1)
	_, err = file.ReadAt(buf, fileLen-1)
	if err != nil {
		return fmt.Errorf("failed to read config file tail: %w", err)
	}
	if buf[0] == 10 {
		s = pwd
	} else {
		s = "\n" + pwd
	}

	_, err = file.Write([]byte(s))
	if err != nil {
		return fmt.Errorf("failed to append to config: %w", err)
	}

	return nil
}

func editMode() (err error) {
	var cmd exec.Cmd
	projectsFilePath := cli.GetProjectsConfig()

	if editor, exists := os.LookupEnv("EDITOR"); exists {
		cmd = *exec.Command(editor, projectsFilePath)
	} else if editor, exists := os.LookupEnv("VISUAL"); exists {
		cmd = *exec.Command(editor, projectsFilePath)
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
	appendMode := flag.Bool("a", false, "Append current directory to projects")
	flag.Parse()

	flags := flags{
		*editMode,
		*appendMode,
	}

	return flags
}
