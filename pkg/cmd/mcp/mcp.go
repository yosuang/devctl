package mcp

import (
	cmdGenconfig "devopsctl/pkg/cmd/mcp/genconfig"

	"github.com/spf13/cobra"
)

func NewCmdMcp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp <command>",
		Short: "Manage MCP (Model Context Protocol) configurations",
		Long:  `Manage MCP (Model Context Protocol) configurations and services for Claude Code integration.`,
	}

	// Add subcommands
	cmd.AddCommand(cmdGenconfig.NewCmdGenConfig())

	return cmd
}
