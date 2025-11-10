package cli

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
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
	var in string
	// var selection int

	for i, entry := range entries {
		fmt.Printf("(%d): %s", i+1, entry)
		if i != len(entries)-1 {
			fmt.Printf("\n")
		}
	}

	if in, err = getUserInput(); err != nil {
		return err
	}

	_, err = strconv.Atoi(in)
	if err != nil {
		return err
	}

	return nil
}

func getUserInput() (s string, err error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to set terminal to raw: %w", err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	var buf [1]byte
	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("failed to read a character from stdin: %w", err)
	}

	return string(buf[0]), nil
}
