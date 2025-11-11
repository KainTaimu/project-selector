package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"golang.org/x/term"
)

type Color int

const (
	Green Color = iota
	Red
)

func RunSelector() (err error) {
	var entries []string
	if entries, err = ReadConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	if err = mainLoop(entries); err != nil {
		return err
	}
	return nil
}

func RunQuickJumper() (err error) {
	var entries []string
	if entries, err = ReadConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	arg := flag.Arg(0)

	selection, err := strconv.Atoi(arg)
	if err != nil {
		return nil
	}

	if selection <= 0 || selection > len(entries) {
		return nil
	}

	if err = startNewBuffer(selection, entries); err != nil {
		return err
	}

	return nil
}

func mainLoop(entries []string) (err error) {
	fmt.Print(SaveCursor)
	printEntries(entries, Green)

	var in string
	if in, err = getUserInput(); err != nil {
		return err
	}

	selection, err := strconv.Atoi(in)
	if err != nil {
		fmt.Print(RestoreCursor + ClearToEnd)
		return nil
	}

	if selection <= 0 || selection > len(entries) {
		fmt.Print(RestoreCursor + ClearToEnd)
		return nil
	}

	if err = startNewBuffer(selection, entries); err != nil {
		return err
	}

	return nil
}

func startNewBuffer(selection int, entries []string) (err error) {
	path := entries[selection-1]
	if path[0] == '~' {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			return fmt.Errorf("$HOME is not set")
		}
		path = filepath.Join(home, path[1:])
		if !IsDir(path) {
			return fmt.Errorf("malformed working dir \"%s\". is $HOME set correctly?", path)
		}
	}

	shell := os.Getenv("SHELL")
	if !IsFile(shell) {
		return fmt.Errorf("$SHELL should be a path to file")
	}

	cmd := exec.Command(shell)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Print(HomeCursor + ClearScreen)
	defer func() {
		fmt.Print(HomeCursor + ClearScreen)
	}()
	if err = cmd.Start(); err != nil {
		return err
	}

	_, _ = cmd.Stdout.Write([]byte(ColorEntry("Entered buffer " + path + "\n")))

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func printEntries(entries []string, color Color) {
	for i, entry := range entries {
		var s string
		switch color {
		case Green:
			s = ColorEntry("(" + strconv.Itoa(i+1) + ")")
		case Red:
			s = ColorPop("(" + strconv.Itoa(i+1) + ")")
		}

		fmt.Printf("%s %s\n", s, entry)
	}
}

func getUserInput() (s string, err error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to set terminal to raw: %w", err)
	}
	fmt.Print(HideCursor)
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print(ShowCursor)
	}()

	var buf [1]byte
	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("failed to read a character from stdin: %w", err)
	}

	return string(buf[0]), nil
}
