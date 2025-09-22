package cmd

import (
	"devopsctl/internal/logging"
	"devopsctl/pkg/cli"
	"devopsctl/pkg/cmdutil"
	"errors"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var settings = cli.New()

func NewCmdRoot() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "devopsctl",
		Short:        "DevOps Management Command Line Interface",
		Long:         `DevOps Management Command Line Interface`,
		Example:      `$ devopsctl mcp list`,
		SilenceUsage: true,
	}

	settings.AddFlags(cmd.PersistentFlags())

	setupLogging()

	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.AddCommand(NewCmdMcp())

	return cmd, nil
}

func setupLogging() {
	logger := logging.NewLogger(func() bool { return settings.Debug })
	slog.SetDefault(logger)
}

func rootFlagErrorFunc(_ *cobra.Command, err error) error {
	if errors.Is(err, pflag.ErrHelp) {
		return err
	}
	return cmdutil.FlagErrorWrap(err)
}

type CommandError struct {
	error
	ExitCode int
}
