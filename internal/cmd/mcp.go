package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdMcp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mcp",
		Short:   "Configure and manage MCP servers",
		Long:    `Configure and manage MCP servers for different clients`,
		Example: `$ devopsctl mcp manage`,
	}

	return cmd
}
