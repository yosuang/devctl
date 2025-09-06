package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ConfigManager handles configuration file operations
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new config manager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
	}
}

// ReadConfig reads and parses the MCP configuration file
func (cm *ConfigManager) ReadConfig() (*MCPConfig, error) {
	file, err := os.Open(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty config if file doesn't exist
			return &MCPConfig{
				MCPServers: make(map[string]MCPServerConfig),
			}, nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config MCPConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize MCPServers if nil
	if config.MCPServers == nil {
		config.MCPServers = make(map[string]MCPServerConfig)
	}

	return &config, nil
}

// WriteConfig writes the MCP configuration to file
func (cm *ConfigManager) WriteConfig(config *MCPConfig) error {
	// Create backup first
	if err := cm.createBackup(); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write to temporary file first
	tempPath := cm.configPath + ".tmp"
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp config file: %w", err)
	}

	// Atomic move
	if err := os.Rename(tempPath, cm.configPath); err != nil {
		os.Remove(tempPath) // Cleanup temp file
		return fmt.Errorf("failed to move temp config file: %w", err)
	}

	return nil
}

// ServerExists checks if a server exists in the configuration
func (cm *ConfigManager) ServerExists(config *MCPConfig, name string) bool {
	_, exists := config.MCPServers[name]
	return exists
}

// AddServer adds a server to the configuration
func (cm *ConfigManager) AddServer(config *MCPConfig, server MCPServer) {
	config.MCPServers[server.Name] = server.ToMCPServerConfig()
}

// RemoveServer removes a server from the configuration
func (cm *ConfigManager) RemoveServer(config *MCPConfig, name string) bool {
	if !cm.ServerExists(config, name) {
		return false
	}
	delete(config.MCPServers, name)
	return true
}

// ListServers returns all servers in the configuration
func (cm *ConfigManager) ListServers(config *MCPConfig) []MCPServer {
	servers := make([]MCPServer, 0, len(config.MCPServers))
	for name, serverConfig := range config.MCPServers {
		servers = append(servers, serverConfig.ToMCPServer(name))
	}
	return servers
}

// createBackup creates a backup of the current config file
func (cm *ConfigManager) createBackup() error {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return nil // No backup needed if file doesn't exist
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := cm.configPath + ".backup." + timestamp

	// Copy file
	src, err := os.Open(cm.configPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
