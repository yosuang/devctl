package devopsctlcmd

import (
	"devopsctl/pkg/cmd/root"
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
	rootCmd, err := root.NewCmdRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create root command: %s\n", err)
		return exitError
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create root command: %s\n", err)
		return exitError
	}

	return exitOK
}
