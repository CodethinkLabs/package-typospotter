package pypi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/codethinklabs/package-typospotter/internal/pkgmgr"
	"github.com/k3a/html2text"
)

type PypiManager struct{}

func New() *PypiManager {
	return &PypiManager{}
}

func (pypi PypiManager) GetName() string {
	return "pypi"
}

func (pypi PypiManager) ListPackages() (pkgmgr.PackageSet, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get("https://pypi.org/simple/")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected non 200 status code %v", resp.StatusCode)
	}
	strBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rawText := html2text.HTML2Text(string(strBody))
	rawText = strings.ReplaceAll(rawText, "/simple/", "")
	rawText = strings.ReplaceAll(rawText, "/", "")
	return strings.Split(rawText, " "), nil
}
