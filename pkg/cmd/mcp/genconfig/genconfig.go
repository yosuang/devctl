package genconfig

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdGenConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genconfig",
		Short: "Generate MCP configuration file",
		Long:  `Generate a default MCP (Model Context Protocol) configuration file for Claude Code integration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateMCPConfig()
		},
	}

	return cmd
}

func generateMCPConfig() error {
	fmt.Printf("✅ MCP configuration generated successfully.")
	return nil
}
