package pkgmgr

type PackageSet []string

type PackageManager interface {
	GetName() string
	ListPackages() (PackageSet, error)
}

func (set PackageSet) Difference(packages PackageSet) []string {
	diff := []string{}
	for _, sa := range set {
		found := false
		for _, sb := range packages {
			if sa == sb {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, sa)
		}
	}
	return diff
}
