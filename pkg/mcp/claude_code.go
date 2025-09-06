package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// ClaudeCodeClient implements MCPClient for Claude Code
type ClaudeCodeClient struct {
	configManager *ConfigManager
}

// Name returns the client name
func (c *ClaudeCodeClient) Name() string {
	return "claude-code"
}

// GetConfigPath returns the config file path for Claude Code
func (c *ClaudeCodeClient) GetConfigPath() (string, error) {
	var configPath string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		configPath = filepath.Join(appData, "Claude", "claude_desktop_config.json")
	case "darwin":
		home := os.Getenv("HOME")
		if home == "" {
			return "", fmt.Errorf("HOME environment variable not set")
		}
		configPath = filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "linux":
		home := os.Getenv("HOME")
		if home == "" {
			return "", fmt.Errorf("HOME environment variable not set")
		}
		configPath = filepath.Join(home, ".config", "claude", "claude_desktop_config.json")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return configPath, nil
}

// getConfigManager returns the config manager, creating it if necessary
func (c *ClaudeCodeClient) getConfigManager() (*ConfigManager, error) {
	if c.configManager == nil {
		configPath, err := c.GetConfigPath()
		if err != nil {
			return nil, err
		}
		c.configManager = NewConfigManager(configPath)
	}
	return c.configManager, nil
}

// ListServers returns all installed MCP servers
func (c *ClaudeCodeClient) ListServers() ([]MCPServer, error) {
	cm, err := c.getConfigManager()
	if err != nil {
		return nil, err
	}

	config, err := cm.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cm.ListServers(config), nil
}

// InstallServer installs an MCP server
func (c *ClaudeCodeClient) InstallServer(server MCPServer) error {
	cm, err := c.getConfigManager()
	if err != nil {
		return err
	}

	config, err := cm.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Check if server already exists
	if cm.ServerExists(config, server.Name) {
		return fmt.Errorf("MCP server '%s' already exists", server.Name)
	}

	// Add server to config
	cm.AddServer(config, server)

	// Write config back
	if err := cm.WriteConfig(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// UninstallServer uninstalls an MCP server
func (c *ClaudeCodeClient) UninstallServer(name string) error {
	cm, err := c.getConfigManager()
	if err != nil {
		return err
	}

	config, err := cm.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Check if server exists
	if !cm.ServerExists(config, name) {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	// Remove server from config
	cm.RemoveServer(config, name)

	// Write config back
	if err := cm.WriteConfig(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
