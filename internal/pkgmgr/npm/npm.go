package npm

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/codethinklabs/package-typospotter/internal/pkgmgr"
	"github.com/tidwall/gjson"
)

type NPMManager struct{}

func New() *NPMManager {
	return &NPMManager{}
}

func (npm NPMManager) GetName() string {
	return "npm"
}

func (npm NPMManager) ListPackages() (pkgmgr.PackageSet, error) {
	packages := pkgmgr.PackageSet{}
	httpClient := &http.Client{}
	resp, err := httpClient.Get("https://replicate.npmjs.com/_all_docs")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected non 200 status code %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := gjson.Get(string(data), "rows.#.id")
	result.ForEach(func(k, v gjson.Result) bool {
		packages = append(packages, v.String())
		return true
	})

	return packages, nil
}
