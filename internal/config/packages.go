package config

func FindPackage(packages []PackageConfig, name string) *PackageConfig {
	for i := range packages {
		if packages[i].Name == name {
			return &packages[i]
		}
	}
	return nil
}

func MergePackages(existing, new []PackageConfig) []PackageConfig {
	pkgMap := make(map[string]PackageConfig)

	for _, pkg := range existing {
		pkgMap[pkg.Name] = pkg
	}

	for _, pkg := range new {
		pkgMap[pkg.Name] = pkg
	}

	result := make([]PackageConfig, 0, len(pkgMap))
	for _, pkg := range pkgMap {
		result = append(result, pkg)
	}

	return result
}
