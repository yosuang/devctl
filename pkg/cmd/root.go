package cmd

import (
	"devopsctl/pkg/cmdutil"
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var verboseMode bool

func NewCmdRoot() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "devopsctl",
		Short:   "DevOps Management Command Line Interface",
		Long:    `DevOps Management Command Line Interface`,
		Example: `$ devopsctl mcp list`,
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if verboseMode {
				configureLogger(slog.LevelDebug)
			}
			return nil
		},
	}

	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.PersistentFlags().BoolVar(&verboseMode, "verbose", false, "Show detailed output")
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.SilenceUsage = true
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(NewCmdMcp())

	return cmd, nil
}

func rootFlagErrorFunc(cmd *cobra.Command, err error) error {
	if errors.Is(err, pflag.ErrHelp) {
		return err
	}
	return cmdutil.FlagErrorWrap(err)
}

func init() {
	configureLogger(slog.LevelWarn)
}

func configureLogger(level slog.Level) {
	opts := &slog.HandlerOptions{Level: level}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}
