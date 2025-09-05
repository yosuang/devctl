package root

import (
	mcpCmd "devopsctl/pkg/cmd/mcp"

	"github.com/spf13/cobra"
)

func NewCmdRoot() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "devopsctl",
		Short:   "Work seamlessly with DevOps from the command line.",
		Long:    `Work seamlessly with DevOps from the command line.`,
		Example: `$ devopsctl mcp genconfig`,
	}

	cmd.AddCommand(mcpCmd.NewCmdMcp())

	return cmd, nil
}
