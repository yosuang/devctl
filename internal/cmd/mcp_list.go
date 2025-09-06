package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"devopsctl/pkg/mcp"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	McpClient  string
	JsonOutput bool
}

func NewCmdList() *cobra.Command {
	o := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed MCP servers",
		Long:  `List all installed MCP servers for the specified client`,
		Example: `  # List MCP servers for Claude Code
  devopsctl mcp list --mcp-client claude-code
  
  # List with JSON output
  devopsctl mcp list --mcp-client claude-code --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run()
		},
	}

	cmd.Flags().StringVar(&o.McpClient, "mcp-client", "", "MCP client to list servers for (required)")
	cmd.Flags().BoolVar(&o.JsonOutput, "json", false, "Output in JSON format")
	cmd.MarkFlagRequired("mcp-client")

	return cmd
}

func (o *ListOptions) run() error {
	// Get the MCP client
	client, err := mcp.GetClient(o.McpClient)
	if err != nil {
		return fmt.Errorf("failed to get MCP client: %w", err)
	}

	// List servers
	servers, err := client.ListServers()
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	// Output results
	if o.JsonOutput {
		return o.outputJSON(servers)
	}

	return o.outputTable(servers)
}

func (o *ListOptions) outputJSON(servers []mcp.MCPServer) error {
	data, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal servers to JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func (o *ListOptions) outputTable(servers []mcp.MCPServer) error {
	if len(servers) == 0 {
		fmt.Printf("No MCP servers found for client '%s'\n", o.McpClient)
		return nil
	}

	// Create tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)

	// Print header
	fmt.Fprintln(w, "NAME\tCOMMAND")
	fmt.Fprintln(w, "----\t-------")

	// Print servers
	for _, server := range servers {
		command := server.Command
		if len(server.Args) > 0 {
			command += " " + strings.Join(server.Args, " ")
		}

		fmt.Fprintf(w, "%s\t%s\n", server.Name, command)
	}

	return w.Flush()
}
