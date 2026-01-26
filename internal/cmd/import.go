package cmd

import (
	"context"
	"devctl/internal/config"
	"devctl/pkg/pkgmgr"
	"devctl/pkg/version"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdImport(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import packages from JSON file",
		Long:  `Import packages from a JSON configuration file and install them using the configured package managers.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runImport(cfg, args[0])
		},
	}
	return cmd
}

func runImport(cfg *config.Config, filePath string) error {
	return runImportWithManager(cfg, filePath, nil)
}

func runImportWithManager(cfg *config.Config, filePath string, testMgr pkgmgr.Manager) error {
	importConfig, err := loadImportConfig(filePath)
	if err != nil {
		return err
	}

	ctx := context.Background()
	successfulPackages := []config.PackageConfig{}

	for _, pkg := range importConfig.Packages {
		if !isValidPackage(pkg) {
			continue
		}

		if pkg.InstalledBy == config.Pwsh {
			if testMgr == nil {
				fmt.Printf("⚠ Skipping pwsh package: %s (pwsh not supported)\n", pkg.Name)
			}
			continue
		}

		if _, ok := cfg.PackageManagers[pkg.InstalledBy]; !ok {
			return fmt.Errorf("package manager %s not configured", pkg.InstalledBy)
		}

		mgr := testMgr
		if mgr == nil {
			mgrConfig := cfg.PackageManagers[pkg.InstalledBy]
			mgr, err = getManager(pkg.InstalledBy, mgrConfig.ExecutablePath)
			if err != nil {
				return fmt.Errorf("failed to get manager: %w", err)
			}
		}

		if err := processPackage(ctx, mgr, pkg); err != nil {
			if testMgr == nil {
				fmt.Printf("✗ Failed to process %s: %v\n", pkg.Name, err)
			}
			continue
		}

		successfulPackages = append(successfulPackages, pkg)
	}

	cfg.Packages = config.MergePackages(cfg.Packages, successfulPackages)

	if err := config.SaveToFile(cfg, cfg.ConfigDir); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	if testMgr == nil {
		fmt.Printf("\n✓ Successfully imported %d package(s)\n", len(successfulPackages))
	}
	return nil
}

func loadImportConfig(filePath string) (*config.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var importConfig config.Config
	if err := json.Unmarshal(data, &importConfig); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &importConfig, nil
}

func isValidPackage(pkg config.PackageConfig) bool {
	return pkg.Name != "" && pkg.Version != "" && pkg.InstalledBy != ""
}

func processPackage(ctx context.Context, mgr pkgmgr.Manager, pkg config.PackageConfig) error {
	installedPackages, err := mgr.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	var installedPkg *pkgmgr.Package
	for i := range installedPackages {
		if installedPackages[i].Name == pkg.Name {
			installedPkg = &installedPackages[i]
			break
		}
	}

	if installedPkg != nil {
		if version.Equal(installedPkg.Version, pkg.Version) {
			return nil
		}

		if err := mgr.Uninstall(ctx, pkg.Name); err != nil {
			return fmt.Errorf("failed to uninstall: %w", err)
		}
	}

	if err := mgr.Install(ctx, pkg.Name); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	return nil
}

func getManager(pm config.PackageManager, execPath string) (pkgmgr.Manager, error) {
	return nil, fmt.Errorf("manager not implemented")
}
