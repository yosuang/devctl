package mcp

import "errors"

// MCPClient represents an interface for different MCP clients
type MCPClient interface {
	Name() string
	GetConfigPath() (string, error)
	ListServers() ([]MCPServer, error)
	InstallServer(server MCPServer) error
	UninstallServer(name string) error
}

// Registry holds all registered MCP clients
var registeredClients = map[string]MCPClient{
	"claude-code": &ClaudeCodeClient{},
}

// GetClient returns a client by name
func GetClient(name string) (MCPClient, error) {
	client, exists := registeredClients[name]
	if !exists {
		return nil, errors.New("unsupported MCP client: " + name)
	}
	return client, nil
}

// ListClients returns all supported client names
func ListClients() []string {
	clients := make([]string, 0, len(registeredClients))
	for name := range registeredClients {
		clients = append(clients, name)
	}
	return clients
}
