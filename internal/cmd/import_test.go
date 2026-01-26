package cmd

import (
	"context"
	"devctl/internal/config"
	"devctl/pkg/pkgmgr"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockManager implements pkgmgr.Manager for testing
type mockManager struct {
	installedPackages []pkgmgr.Package
	installErr        error
	uninstallErr      error
	listErr           error
	installedCalls    []string
	uninstalledCalls  []string
}

func (m *mockManager) Install(ctx context.Context, names ...string) error {
	m.installedCalls = append(m.installedCalls, names...)
	return m.installErr
}

func (m *mockManager) Uninstall(ctx context.Context, names ...string) error {
	m.uninstalledCalls = append(m.uninstalledCalls, names...)
	return m.uninstallErr
}

func (m *mockManager) List(ctx context.Context) ([]pkgmgr.Package, error) {
	return m.installedPackages, m.listErr
}

func TestImport(t *testing.T) {
	tests := []struct {
		name              string
		jsonContent       string
		existingPackages  []config.PackageConfig
		installedPackages []pkgmgr.Package
		installErr        error
		uninstallErr      error
		listErr           error
		wantErr           bool
		wantInstalled     []string
		wantUninstalled   []string
		wantPackages      []string
	}{
		{
			name: "valid JSON with new package",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.40.0", "installedBy": "scoop"}
				]
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			wantErr:           false,
			wantInstalled:     []string{"git"},
			wantUninstalled:   []string{},
			wantPackages:      []string{"git"},
		},
		{
			name: "invalid JSON format",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.40.0"
				]
			}`,
			existingPackages: []config.PackageConfig{},
			wantErr:          true,
		},
		{
			name: "missing required field - name",
			jsonContent: `{
				"packages": [
					{"version": "2.40.0", "installedBy": "scoop"}
				]
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			wantErr:           false,
			wantInstalled:     []string{},
			wantUninstalled:   []string{},
			wantPackages:      []string{},
		},
		{
			name: "missing required field - version",
			jsonContent: `{
				"packages": [
					{"name": "git", "installedBy": "scoop"}
				]
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			wantErr:           false,
			wantInstalled:     []string{},
			wantUninstalled:   []string{},
			wantPackages:      []string{},
		},
		{
			name: "missing required field - installedBy",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.40.0"}
				]
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			wantErr:           false,
			wantInstalled:     []string{},
			wantUninstalled:   []string{},
			wantPackages:      []string{},
		},
		{
			name: "already installed with matching version",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.40.0", "installedBy": "scoop"}
				]
			}`,
			existingPackages: []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{
				{Name: "git", Version: "2.40.0"},
			},
			wantErr:         false,
			wantInstalled:   []string{},
			wantUninstalled: []string{},
			wantPackages:    []string{"git"},
		},
		{
			name: "already installed with different version",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.41.0", "installedBy": "scoop"}
				]
			}`,
			existingPackages: []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{
				{Name: "git", Version: "2.40.0"},
			},
			wantErr:         false,
			wantInstalled:   []string{"git"},
			wantUninstalled: []string{"git"},
			wantPackages:    []string{"git"},
		},
		{
			name: "pwsh package should be skipped",
			jsonContent: `{
				"packages": [
					{"name": "powershell", "version": "7.3.0", "installedBy": "pwsh"}
				]
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			wantErr:           false,
			wantInstalled:     []string{},
			wantUninstalled:   []string{},
			wantPackages:      []string{},
		},
		{
			name: "install failure should continue processing",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.40.0", "installedBy": "scoop"},
					{"name": "vim", "version": "9.0.0", "installedBy": "scoop"}
				]
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			installErr:        pkgmgr.ErrNotFound,
			wantErr:           false,
			wantInstalled:     []string{"git", "vim"},
			wantUninstalled:   []string{},
			wantPackages:      []string{},
		},
		{
			name: "empty packages array",
			jsonContent: `{
				"packages": []
			}`,
			existingPackages:  []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{},
			wantErr:           false,
			wantInstalled:     []string{},
			wantUninstalled:   []string{},
			wantPackages:      []string{},
		},
		{
			name: "multiple packages with mixed scenarios",
			jsonContent: `{
				"packages": [
					{"name": "git", "version": "2.40.0", "installedBy": "scoop"},
					{"name": "vim", "version": "9.0.0", "installedBy": "scoop"},
					{"name": "curl", "version": "8.0.0", "installedBy": "scoop"}
				]
			}`,
			existingPackages: []config.PackageConfig{},
			installedPackages: []pkgmgr.Package{
				{Name: "git", Version: "2.40.0"},
				{Name: "vim", Version: "8.0.0"},
			},
			wantErr:         false,
			wantInstalled:   []string{"vim", "curl"},
			wantUninstalled: []string{"vim"},
			wantPackages:    []string{"git", "vim", "curl"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// #given: temporary directory and JSON file
			tmpDir := t.TempDir()
			jsonFile := filepath.Join(tmpDir, "import.json")
			err := os.WriteFile(jsonFile, []byte(tt.jsonContent), 0644)
			require.NoError(t, err)

			// #given: config with mock manager
			cfg := &config.Config{
				ConfigDir: tmpDir,
				Packages:  tt.existingPackages,
				PackageManagers: map[config.PackageManager]config.PackageManagerConfig{
					config.Scoop: {
						ID:             config.Scoop,
						ExecutablePath: "/usr/bin/scoop",
					},
				},
			}

			mock := &mockManager{
				installedPackages: tt.installedPackages,
				installErr:        tt.installErr,
				uninstallErr:      tt.uninstallErr,
				listErr:           tt.listErr,
			}

			// #when: runImport is called
			err = runImportWithManager(cfg, jsonFile, mock)

			// #then: verify error expectation
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// #then: verify install/uninstall calls
			require.ElementsMatch(t, tt.wantInstalled, mock.installedCalls)
			require.ElementsMatch(t, tt.wantUninstalled, mock.uninstalledCalls)

			// #then: verify config was updated correctly
			if len(tt.wantPackages) > 0 {
				savedConfig := loadConfigFromFile(t, tmpDir)
				var packageNames []string
				for _, pkg := range savedConfig.Packages {
					packageNames = append(packageNames, pkg.Name)
				}
				require.ElementsMatch(t, tt.wantPackages, packageNames)
			}
		})
	}
}

func TestNewCmdImport(t *testing.T) {
	// #given: a config
	cfg := &config.Config{
		ConfigDir: t.TempDir(),
	}

	// #when: NewCmdImport is called
	cmd := NewCmdImport(cfg)

	// #then: command should be properly configured
	require.NotNil(t, cmd)
	require.Equal(t, "import", cmd.Use)
	require.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)
}

// Helper function to load config from file
func loadConfigFromFile(t *testing.T, configDir string) *config.Config {
	configPath := filepath.Join(configDir, config.AppName+".json")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var cfg config.Config
	err = json.Unmarshal(data, &cfg)
	require.NoError(t, err)

	return &cfg
}
