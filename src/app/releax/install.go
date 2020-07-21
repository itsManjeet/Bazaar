package releax

import (
	"fmt"
	"log"
	"path"

	"github.com/itsmanjeet/bazaar/src/app/store"
)

// Install install app from releax store
func (s Store) Install(appname string, progfunc store.ProgFunc) error {

	log.Println("Processing Dependencies for", appname)
	apd, err := s.Depends(appname)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, a := range apd {
		if err := s.Install(a.Name(), progfunc); err != nil {
			return err
		}
	}
	log.Println("Installing", appname)

	app, err := s.GetApp(appname)
	if err != nil {
		log.Println("Error while get app", appname)
		return err
	}

	url := fmt.Sprintf("%s/%s", s.BinaryURL, app.PackageFile())
	log.Println("downloading ", url)
	downloader := store.Downloader{
		ProgressFunction: progfunc,
	}

	pkgfile := path.Join(s.DataDir, app.PackageFile())
	if err := downloader.Download(pkgfile, url); err != nil {
		log.Println("error while downloading file", err)
		return err
	}

	log.Println("Extracting ", pkgfile)
	if err := store.Extractor(pkgfile, "/"); err != nil {
		log.Println("error while extraction File", err)
		return err
	}

	return nil
}
