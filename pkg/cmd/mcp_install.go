package cmd

import (
	"devopsctl/pkg/cmdutil"
	"log/slog"

	"github.com/spf13/cobra"
)

type InstallOptions struct {
	ServerName string
	McpClient  string
	ExtraArgs  []string
}

func newCmdInstall() *cobra.Command {
	o := &InstallOptions{}

	cmd := &cobra.Command{
		Use:     "install <name>",
		Short:   "Install MCP server",
		Long:    `Install MCP server`,
		Example: `$ devopsctl mcp install context7 --client claude-code -- cmd /c npx -y @upstash/context7-mcp`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.ServerName = args[0]
			o.ExtraArgs = args[1:]

			return o.run()
		},
	}

	cmd.Flags().StringVarP(&o.McpClient, "client", "c", "", "MCP client to install server for (e.g. claude-code)")

	return cmd
}

func (o *InstallOptions) run() error {
	if o.McpClient == "" {
		return cmdutil.FlagErrorf("--client flag is required")
	}

	if len(o.ExtraArgs) > 0 {
		slog.Debug("", "ExtraArgs", o.ExtraArgs)
	}

	// TODO: Implement actual MCP server installation logic

	return nil
}
