package platform

import (
	"devctl/pkg/pkgmgr"
	"runtime"
)

// Platform represents an operating system platform.
type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformDarwin  Platform = "darwin"
	PlatformLinux   Platform = "linux"
)

// GetCurrent returns the current platform.
func GetCurrent() Platform {
	return Platform(runtime.GOOS)
}

// GetSupportedManagers returns the list of package managers supported on the given platform.
func GetSupportedManagers(p Platform) []pkgmgr.ManagerType {
	switch p {
	case PlatformWindows:
		return []pkgmgr.ManagerType{
			pkgmgr.ManagerTypeScoop,
			pkgmgr.ManagerTypePwsh,
		}
	case PlatformDarwin:
		return []pkgmgr.ManagerType{
			pkgmgr.ManagerTypeBrew,
		}
	case PlatformLinux:
		return []pkgmgr.ManagerType{
			pkgmgr.ManagerTypeBrew,
			pkgmgr.ManagerTypeApt,
		}
	default:
		return []pkgmgr.ManagerType{}
	}
}

// IsManagerSupported checks if a package manager is supported on the given platform.
func IsManagerSupported(mgr pkgmgr.ManagerType, p Platform) bool {
	supported := GetSupportedManagers(p)
	for _, m := range supported {
		if m == mgr {
			return true
		}
	}
	return false
}
