package flatpak

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/itsmanjeet/bazaar/src/app/store"
)

// Store contain flatpak apps
type Store struct {
	AppStream string
	Cache     []store.App
}

func (s Store) ID() string {
	return "flatpak"
}

// ListApps return list of flatpak apps
func (s Store) ListApps() ([]store.App, error) {
	apfptr, err := ioutil.ReadFile(s.AppStream)
	if err != nil {
		return nil, err
	}

	var c AppStreamData
	err = xml.Unmarshal(apfptr, &c)
	if err != nil {
		return nil, err
	}

	aplst := make([]store.App, 0)
	for _, a := range c.Components {
		aplst = append(aplst, a)
	}

	return aplst, nil
}

func (s Store) GetApp(appname string) (store.App, error) {
	apl, err := s.ListApps()
	if err != nil {
		return nil, err
	}
	for _, a := range apl {
		if a.Name() == appname {
			return a, nil
		}
	}
	return nil, errors.New("not in repo")
}

func (s Store) IsInstalled(name string) bool {
	list, err := exec.Command("flatpak", "list", "--columns=app").Output()
	if err != nil {
		return false
	}

	appsList := strings.Split(string(list), "\n")

	app, err := s.GetApp(name)
	if err != nil {
		return false
	}

	appID := strings.ReplaceAll(app.ID(), ".desktop", "")
	for _, a := range appsList {
		if appID == a {
			return true
		}
	}

	return false
}

func (s Store) Depends(appname string) ([]store.App, error) {
	return nil, nil
}
