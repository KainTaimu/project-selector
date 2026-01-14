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

// RunSelector runs the interactive menu to select the directory to jump to
func RunSelector() (err error) {
	var entries []Entry
	if entries, err = ReadConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	if err = mainLoop(entries); err != nil {
		return err
	}
	return nil
}

// RunQuickJumper takes the first argument as the index of the directory to jump to as
// displayed by RunSelector
func RunQuickJumper() (err error) {
	fmt.Print(SaveCursor)
	var entries []Entry
	if entries, err = ReadConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	arg := flag.Arg(0)

	selection, err := strconv.Atoi(arg)
	if err != nil {
		return nil
	}

	var entry *Entry
	if entry, err = checkSelectionIsValid(selection, entries); err != nil {
		return err
	}

	if err = startNewBuffer(entry.Path); err != nil {
		return err
	}

	return nil
}

func mainLoop(entries []Entry) (err error) {
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

	var entry *Entry
	if entry, err = checkSelectionIsValid(selection, entries); err != nil {
		return err
	}

	if err = startNewBuffer(entry.Path); err != nil {
		return err
	}

	return nil
}

func checkSelectionIsValid(selection int, entries []Entry) (*Entry, error) {
	if selection <= 0 || selection > len(entries) {
		fmt.Print(RestoreCursor + ClearToEnd)
		return nil, fmt.Errorf("selection out of bounds")
	}

	entry := &entries[selection-1]
	if !entry.IsValid {
		fmt.Print(RestoreCursor + ClearToEnd)
		return nil, fmt.Errorf("entry \"%s\" is invalid", entry.Path)
	}
	return entry, nil
}

func startNewBuffer(path string) (err error) {
	path = TildeExpansion(path)

	shell := os.Getenv("SHELL")
	if !IsFile(shell) {
		return fmt.Errorf("$SHELL should be a path to the shell program")
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

	_, _ = cmd.Stdout.Write([]byte(ColorEntry("Entered buffer " + ShortenTildeExpansion(path) + "\n")))

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func printEntries(entries []Entry, color Color) {
	for i, entry := range entries {
		var s string
		switch color {
		case Green:
			s = ColorEntry("(" + strconv.Itoa(i+1) + ")")
		case Red:
			s = ColorPop("(" + strconv.Itoa(i+1) + ")")
		}

		path := entry.Path
		path = ShortenTildeExpansion(path)
		path = filepath.Clean(path)

		if !entry.IsValid {
			s = ColorInvalid("(" + strconv.Itoa(i+1) + ")")
			path = ColorInvalid(path)
		}

		fmt.Printf("%s %s\n", s, path)
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
