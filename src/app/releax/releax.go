package releax

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/itsmanjeet/bazaar/src/app/store"
)

// Store contain releax apps
type Store struct {
	RecipiePath  string
	SourceURL    string
	BinaryURL    string
	Repositories []string
	Cache        []store.App
	DataDir      string
}

func (s Store) ID() string {
	return "releax"
}

// ListApps return list of apps
func (s Store) ListApps() ([]store.App, error) {
	apl := make([]store.App, 0)

	for _, i := range s.Repositories {
		repopath := path.Join(s.RecipiePath, i)
		if exists(repopath) {
			repo, err := ioutil.ReadDir(repopath)
			if err != nil {
				log.Println("error while loading ", repopath)
				continue
			}
			for _, a := range repo {
				if !a.IsDir() {
					continue
				}
				rpc := path.Join(repopath, a.Name(), "recipie")
				if !exists(rpc) {
					log.Println("warning, recipie not exist for", a.Name(), "in", rpc)
					continue
				}

				recipie, err := NewFrom(path.Join(repopath, a.Name()))
				if err != nil {
					log.Println("error while reading recipie file", rpc, err)
					continue
				}
				apl = append(apl, recipie)
			}
		}
	}

	return apl, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetApp from cache
func (s Store) GetApp(name string) (store.App, error) {
	apl, err := s.ListApps()
	if err != nil {
		return nil, err
	}
	for _, a := range apl {
		if a.Name() == name {
			return a, nil
		}
	}
	return nil, errors.New("not in repo")
}

// IsInstalled check if app is installed
func (s Store) IsInstalled(name string) bool {
	return exists(path.Join(s.DataDir, name, "info"))
}

// Depends return dependencies
func (s Store) Depends(appname string) ([]store.App, error) {
	apl := make([]store.App, 0)
	a, err := s.GetApp(appname)
	if err != nil {
		return nil, err
	}
	for _, p := range a.Depends() {
		app, err := s.GetApp(p)
		if err != nil {
			continue
		}
		if s.IsInstalled(app.Name()) || cacheContain(apl, app.Name()) {
			continue
		} else {
			subdep, err := s.Depends(app.Name())
			if err != nil {
				continue
			}
			for _, a := range subdep {
				if s.IsInstalled(a.Name()) || cacheContain(apl, a.Name()) {
					continue
				}
				apl = append(apl, a)
			}
			apl = append(apl, app)
		}
	}

	return apl, nil
}

func cacheContain(apl []store.App, x string) bool {
	for _, a := range apl {
		if a.Name() == x {
			return true
		}
	}
	return false
}
