package cmd

import (
	"fmt"

	"devopsctl/pkg/mcp"

	"github.com/spf13/cobra"
)

type UninstallOptions struct {
	ServerName string
	McpClient  string
}

func NewCmdUninstall() *cobra.Command {
	o := &UninstallOptions{}

	cmd := &cobra.Command{
		Use:   "uninstall <name>",
		Short: "Uninstall MCP server",
		Long:  `Uninstall an MCP server from the specified client configuration`,
		Example: `  # Uninstall context7 MCP server from Claude Code
  devopsctl mcp uninstall context7 --mcp-client claude-code`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.ServerName = args[0]
			return o.run()
		},
	}

	cmd.Flags().StringVar(&o.McpClient, "mcp-client", "", "MCP client to uninstall server from (required)")
	cmd.MarkFlagRequired("mcp-client")

	return cmd
}

func (o *UninstallOptions) run() error {
	// Get the MCP client
	client, err := mcp.GetClient(o.McpClient)
	if err != nil {
		return fmt.Errorf("failed to get MCP client: %w", err)
	}

	// Uninstall server
	if err := client.UninstallServer(o.ServerName); err != nil {
		return fmt.Errorf("failed to uninstall server: %w", err)
	}

	fmt.Printf("Successfully uninstalled MCP server '%s' from %s\n", o.ServerName, o.McpClient)
	return nil
}
