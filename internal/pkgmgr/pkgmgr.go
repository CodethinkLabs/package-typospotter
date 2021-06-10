package pkgmgr

type PackageSet []string

type PackageManager interface {
	GetName() string
	ListPackages() (PackageSet, error)
}

func (set PackageSet) ListNewPackages(latestPkgSet PackageSet) []string {
	newPkgs := []string{}
	oldPkgMap := map[string]int{}
	for i, pkg := range set {
		oldPkgMap[pkg] = i
	}
	for _, pkg := range latestPkgSet {
		if _, ok := oldPkgMap[pkg]; !ok {
			newPkgs = append(newPkgs, pkg)
		}
	}
	return newPkgs
}
