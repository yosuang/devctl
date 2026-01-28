package cmd

import (
	"context"
	"devctl/internal/config"
	"devctl/internal/installer"
	"devctl/internal/ui"
	"devctl/pkg/executil"
	"devctl/pkg/pkgmgr"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func NewCmdInit(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration by detecting package managers",
		Long:  `Detects installed package managers and saves their information to the configuration file.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runInit(cfg)
		},
	}

	return cmd
}

func runInit(cfg *config.Config) error {
	currentPlatform := pkgmgr.GetCurrent()
	detectResult := detectPackageManagers(currentPlatform)
	displayDetectionResults(detectResult, currentPlatform)

	uninstalled := getUninstalledManagers(detectResult)
	if len(uninstalled) == 0 {
		return saveConfiguration(cfg, detectResult)
	}

	fmt.Println()
	confirmed, err := ui.ConfirmAutoInstall(len(uninstalled))
	if err != nil {
		return fmt.Errorf("failed to get user confirmation: %w", err)
	}

	if !confirmed {
		fmt.Println("\nManual installation guides:")
		for _, mgr := range uninstalled {
			showManualInstallGuide(mgr.Type, string(currentPlatform))
		}
		return saveConfiguration(cfg, detectResult)
	}

	fmt.Println()
	for _, mgr := range uninstalled {
		if err := attemptAutoInstall(mgr.Type, string(currentPlatform)); err != nil {
			fmt.Printf("✗ Failed to install %s: %v\n", mgr.Type, err)
			continue
		}

		path := executil.LookPath(string(mgr.Type))
		if path != "" {
			detectResult[mgr.Type] = PackageManagerInfo{
				Type:           mgr.Type,
				Installed:      true,
				ExecutablePath: path,
			}
		}
	}

	return saveConfiguration(cfg, detectResult)
}

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	failStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	titleStyle   = lipgloss.NewStyle().Bold(true)
)

type PackageManagerInfo struct {
	Type           pkgmgr.ManagerType
	Installed      bool
	ExecutablePath string
}

func detectPackageManagers(p pkgmgr.Platform) map[pkgmgr.ManagerType]PackageManagerInfo {
	managers := map[pkgmgr.ManagerType]PackageManagerInfo{}
	supportedManagers := pkgmgr.GetSupportedManagers(p)

	for _, mgr := range supportedManagers {
		path := executil.LookPath(string(mgr))

		managers[mgr] = PackageManagerInfo{
			Type:           mgr,
			Installed:      path != "",
			ExecutablePath: path,
		}
	}

	return managers
}

func displayDetectionResults(results map[pkgmgr.ManagerType]PackageManagerInfo, p pkgmgr.Platform) {
	fmt.Println(titleStyle.Render(fmt.Sprintf("\nPackage Manager Detection (%s)", p)))
	fmt.Println(strings.Repeat("─", 50))

	for _, mgr := range results {
		if mgr.Installed {
			fmt.Printf("%s %-10s Installed at: %s\n",
				successStyle.Render("✓"),
				mgr.Type,
				mgr.ExecutablePath)
		} else {
			fmt.Printf("%s %-10s Not installed\n",
				failStyle.Render("✗"),
				mgr.Type)
		}
	}
}

func getUninstalledManagers(results map[pkgmgr.ManagerType]PackageManagerInfo) []PackageManagerInfo {
	var uninstalled []PackageManagerInfo
	for _, mgr := range results {
		if !mgr.Installed {
			uninstalled = append(uninstalled, mgr)
		}
	}
	return uninstalled
}

func saveConfiguration(cfg *config.Config, results map[pkgmgr.ManagerType]PackageManagerInfo) error {
	packageManagers := map[pkgmgr.ManagerType]config.PackageManagerConfig{}
	for _, p := range results {
		if p.Installed {
			packageManagers[p.Type] = config.PackageManagerConfig{
				Type:           p.Type,
				Version:        "",
				ExecutablePath: p.ExecutablePath,
			}
		}
	}
	cfg.PackageManagers = packageManagers

	if err := config.SaveToFile(cfg, cfg.ConfigDir); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	configPath := fmt.Sprintf("%s/%s.json", cfg.ConfigDir, config.AppName)
	fmt.Printf("\n%s %s\n", successStyle.Render("✓"), fmt.Sprintf("Configuration saved to: %s", configPath))

	return nil
}

func attemptAutoInstall(managerType pkgmgr.ManagerType, platformStr string) error {
	inst := installer.GetInstaller(managerType)
	if inst == nil {
		return fmt.Errorf("no installer available for %s", managerType)
	}

	canAuto, err := inst.CanAutoInstall()
	if !canAuto {
		fmt.Printf("\n%s %s: Automatic installation not available\n",
			failStyle.Render("✗"),
			managerType)

		showGuide, _ := ui.ConfirmShowGuide()
		if showGuide {
			showManualInstallGuide(managerType, platformStr)
		}
		return fmt.Errorf("automatic installation not supported: %w", err)
	}

	fmt.Printf("\n%s Installing %s...\n", infoStyle.Render("→"), managerType)
	fmt.Println(strings.Repeat("─", 50))

	prereqs := inst.GetPrerequisites()
	fmt.Println("Prerequisites:")
	for _, prereq := range prereqs {
		status := successStyle.Render("✓")
		if !prereq.Passed {
			status = failStyle.Render("✗")
		}
		fmt.Printf("  %s %s: %s\n", status, prereq.Name, prereq.Message)
	}

	cmd := inst.GetInstallCommand()
	fmt.Printf("\nCommand to execute:\n  %s\n\n", infoStyle.Render(cmd))

	confirmed, err := ui.ConfirmProceed(string(managerType))
	if err != nil || !confirmed {
		return fmt.Errorf("installation cancelled by user")
	}

	progressChan := make(chan installer.InstallProgress, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- inst.Install(ctx, progressChan)
		close(progressChan)
	}()

	for progress := range progressChan {
		fmt.Printf("%s %s\n", infoStyle.Render("→"), progress.Message)
	}

	if err := <-errChan; err != nil {
		fmt.Printf("\n%s Installation failed\n", failStyle.Render("✗"))
		showManualInstallGuide(managerType, platformStr)
		return err
	}

	fmt.Printf("\n%s %s installed successfully!\n", successStyle.Render("✓"), managerType)
	return nil
}

func showManualInstallGuide(managerType pkgmgr.ManagerType, platformStr string) {
	guide := installer.GetInstallGuide(managerType, platformStr)
	if guide == nil {
		fmt.Printf("No installation guide available for %s\n", managerType)
		return
	}

	fmt.Printf("\n%s Manual Installation Guide for %s\n", titleStyle.Render("📖"), managerType)
	fmt.Println(strings.Repeat("─", 50))

	for i, instruction := range guide.Instructions {
		fmt.Printf("%d. %s\n", i+1, instruction)
	}

	if guide.URL != "" {
		fmt.Printf("\nMore info: %s\n", infoStyle.Render(guide.URL))
	}

	if guide.VerifyCmd != "" {
		fmt.Printf("Verify installation: %s\n", infoStyle.Render(guide.VerifyCmd))
	}

	fmt.Println()
}
