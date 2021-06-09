package typospotter

import (
	"fmt"
	"sync"

	"github.com/agnivade/levenshtein"
	"github.com/codethinklabs/package-typospotter/internal/pkgmgr"
)

type TypoSpotter struct {
	distanceThreshold int

	pkgMgrs      []pkgmgr.PackageManager
	previousPkgs map[pkgmgr.PackageManager]pkgmgr.PackageSet

	mapMutex sync.Mutex
}

type PotentialSquatter struct {
	NewPkg      string
	ExistingPkg string

	PackageManagerName string
}

type pollMgrResult struct {
	potentialSquatters []PotentialSquatter
	err                error

	pkgMgrName string
}

func New(pkgMgrs []pkgmgr.PackageManager, distanceThreshold int) *TypoSpotter {
	spotter := &TypoSpotter{
		pkgMgrs:           pkgMgrs,
		distanceThreshold: distanceThreshold,
		previousPkgs:      map[pkgmgr.PackageManager]pkgmgr.PackageSet{},
	}
	spotter.PollAndCheck()
	return spotter
}

func (spotter *TypoSpotter) PollAndCheck() (map[string][]PotentialSquatter, []error) {
	resultChan := make(chan pollMgrResult, len(spotter.pkgMgrs))
	errs := []error{}

	suspects := map[string][]PotentialSquatter{}
	for _, mgr := range spotter.pkgMgrs {
		go spotter.pollMgr(mgr, resultChan)
	}
	for range spotter.pkgMgrs {
		result := <-resultChan
		if result.err != nil {
			errs = append(errs, result.err)
			continue
		}
		suspects[result.pkgMgrName] = result.potentialSquatters
	}
	return suspects, errs
}

func (spotter *TypoSpotter) pollMgr(mgr pkgmgr.PackageManager, resultChan chan<- pollMgrResult) {
	suspectSquatters := []PotentialSquatter{}

	pkgs, err := mgr.ListPackages()
	if err != nil {
		resultChan <- pollMgrResult{pkgMgrName: mgr.GetName(), err: err}
		return
	}

	lastPkgs, ok := spotter.previousPkgs[mgr]
	if ok {
		diff := pkgs.Difference(lastPkgs)
		fmt.Printf("Polling complete, found %v packages in %s, %v are new packages\n", len(pkgs), mgr.GetName(), len(diff))
		for _, newPkg := range diff {
			for _, pkg := range lastPkgs {
				distance := levenshtein.ComputeDistance(newPkg, pkg)
				if distance <= spotter.distanceThreshold {
					suspectSquatters = append(suspectSquatters, PotentialSquatter{
						NewPkg:      newPkg,
						ExistingPkg: pkg,
					})
				}
			}
		}
	}
	spotter.mapMutex.Lock()
	spotter.previousPkgs[mgr] = pkgs
	spotter.mapMutex.Unlock()
	resultChan <- pollMgrResult{potentialSquatters: suspectSquatters, pkgMgrName: mgr.GetName()}
}
