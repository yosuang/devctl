package mcp

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Name    string            `json:"name"`
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// MCPConfig represents the MCP configuration structure for Claude Code
type MCPConfig struct {
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
}

// MCPServerConfig represents the server configuration in Claude Code format
type MCPServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// ToMCPServer converts MCPServerConfig to MCPServer
func (c MCPServerConfig) ToMCPServer(name string) MCPServer {
	return MCPServer{
		Name:    name,
		Command: c.Command,
		Args:    c.Args,
		Env:     c.Env,
	}
}

// ToMCPServerConfig converts MCPServer to MCPServerConfig
func (s MCPServer) ToMCPServerConfig() MCPServerConfig {
	return MCPServerConfig{
		Command: s.Command,
		Args:    s.Args,
		Env:     s.Env,
	}
}
