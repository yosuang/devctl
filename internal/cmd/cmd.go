package cmd

import (
	"devopsctl/pkg/cmdutil"
	"errors"
	"fmt"
	"os"
)

type ExitCode int

const (
	exitOK      ExitCode = 0
	exitError   ExitCode = 1
	exitCancel  ExitCode = 2
	exitAuth    ExitCode = 4
	exitPending ExitCode = 8
)

func Main() ExitCode {
	rootCmd, err := NewCmdRoot()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create root command.", "err", err)
		return exitError
	}

	if err := rootCmd.Execute(); err != nil {

		if errors.Is(cmdutil.SilentError, err) {
			return exitError
		}

		return exitError
	}

	return exitOK
}
