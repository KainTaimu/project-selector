package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"project_selector/src/cli"
)

type flags struct {
	Edit bool
}

func main() {
	flags := parseFlags()

	if flags.Edit {
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

func editMode() (err error) {
	var cmd exec.Cmd
	projectsFilePath := os.Getenv(cli.ConfigHomeEnv) + "/" + cli.AppConfigDir + cli.ProjectEntriesFile

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

func parseFlags() (flags flags) {
	edit := flag.Bool("e", false, "Launch an editor set by $EDITOR or $VISUAL")
	flag.Parse()

	if *edit {
		flags.Edit = true
	}

	return flags
}
