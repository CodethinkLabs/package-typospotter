package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/codethinklabs/package-typospotter/internal/pkgmgr"
	"github.com/codethinklabs/package-typospotter/internal/pkgmgr/npm"
	"github.com/codethinklabs/package-typospotter/internal/pkgmgr/pypi"
	"github.com/codethinklabs/package-typospotter/internal/typospotter"

	_ "net/http/pprof"
)

func hostPprof() {
	runtime.SetCPUProfileRate(100)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func main() {
	hostPprof()

	spotter := typospotter.New([]pkgmgr.PackageManager{
		*pypi.New(),
		*npm.New(),
	}, 1)

	for {
		fmt.Println("Polling for potential typosquatters...")
		pkgMgrSuspects, errs := spotter.PollAndCheck()
		for _, err := range errs {
			fmt.Printf("Encountered error: %v\n", err)
		}
		for mgrName, suspects := range pkgMgrSuspects {
			if len(suspects) == 0 {
				fmt.Printf("No suspect packages identified whilst polling %s\n", mgrName)
			}
			for _, suspect := range suspects {
				fmt.Printf("New package %s has been detected as a potential typosquatter for existing package %s in %s\n", suspect.NewPkg, suspect.ExistingPkg, suspect.PackageManagerName)
			}
		}

		time.Sleep(time.Minute * 5)
	}
}
