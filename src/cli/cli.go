package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/term"
)

const (
	clearScreen = "\033[J"
	shell       = "/bin/fish"
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

func mainLoop(entries []string) (err error) {
	printEntries(entries, Green)

	var in string
	if in, err = getUserInput(); err != nil {
		return err
	}

	selection, err := strconv.Atoi(in)
	if err != nil {
		return nil
	}

	if selection <= 0 || selection > len(entries) {
		return nil
	}

	path := entries[selection-1]

	cmd := exec.Command(shell, "-C", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Print("\033[H" + clearScreen)
	defer func() {
		fmt.Print("\033[H" + clearScreen)
	}()

	if err = cmd.Run(); err != nil {
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

		fmt.Printf("%s %s", s, entry)
		if i != len(entries)-1 {
			fmt.Printf("\n")
		}
	}
}

func getUserInput() (s string, err error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to set terminal to raw: %w", err)
	}
	fmt.Print("\033[?25l")
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\033[?25h")
	}()

	var buf [1]byte
	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("failed to read a character from stdin: %w", err)
	}

	return string(buf[0]), nil
}
