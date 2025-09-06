package cmd

import (
	"devopsctl/pkg/mcp"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type InstallOptions struct {
	ServerName string
	McpClient  string
	ExtraArgs  []string
}

func NewCmdInstall() *cobra.Command {
	o := &InstallOptions{}

	cmd := &cobra.Command{
		Use:   "install <name>",
		Short: "Install MCP server",
		Long:  `Install an MCP server to the specified client configuration`,
		Example: `  # Install context7 MCP server to Claude Code
  devopsctl mcp install context7 --mcp-client claude-code -- cmd /c npx -y @upstash/context7-mcp
  
  # Install a Python MCP server
  devopsctl mcp install myserver --mcp-client claude-code -- python /path/to/server.py`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.ServerName = args[0]
			o.ExtraArgs = args[1:]

			return o.run()
		},
	}

	cmd.Flags().StringVar(&o.McpClient, "mcp-client", "", "MCP client to install server for (required)")
	cmd.MarkFlagRequired("mcp-client")

	return cmd
}

func (o *InstallOptions) run() error {
	// Parse command and args from ExtraArgs
	if len(o.ExtraArgs) == 0 {
		return fmt.Errorf("command is required after '--' separator")
	}

	command := o.ExtraArgs[0]
	args := o.ExtraArgs[1:]

	// Get the MCP client
	client, err := mcp.GetClient(o.McpClient)
	if err != nil {
		return fmt.Errorf("failed to get MCP client: %w", err)
	}

	// Create server configuration
	server := mcp.MCPServer{
		Name:    o.ServerName,
		Command: command,
		Args:    args,
		Env:     make(map[string]string),
	}

	// Install server
	if err := client.InstallServer(server); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("MCP server '%s' already exists. Use 'devopsctl mcp uninstall %s --mcp-client %s' to remove it first", o.ServerName, o.ServerName, o.McpClient)
		}
		return fmt.Errorf("failed to install server: %w", err)
	}

	fmt.Printf("Successfully installed MCP server '%s' to %s\n", o.ServerName, o.McpClient)
	return nil
}
