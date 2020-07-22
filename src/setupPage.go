package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/itsmanjeet/bazaar/src/app/store"
)

func setupPage(app store.App) {
	if app == nil {
		log.Println("App is nul", app)
		return
	}
	appImage := getWidget("logoImage").(*gtk.Image)
	appName := getWidget("nameLabel").(*gtk.Label)

	getWidget("descLabel").(*gtk.Label).SetText(app.Description())
	getWidget("licenseLabel").(*gtk.Label).SetText(app.License())
	getWidget("urlLabel").(*gtk.Label).SetText(app.Url())

	buttonBox := getWidget("buttonBox").(*gtk.Box)

	imagePix := app.Icon(128)
	if imagePix != nil {
		appImage.SetFromPixbuf(imagePix)
	}

	if app.Store() != "releax" {
		appName.SetText(fmt.Sprintf("%s %s (by %s)", app.Name(), app.Version(), app.Store()))
	} else {
		appName.SetText(fmt.Sprintf("%s %s-%s", app.Name(), app.Version(), app.Release()))
	}

	isInstalled := false
	for _, a := range stores {
		if a.IsInstalled(app.Name()) {
			isInstalled = true
			break
		}
	}

	if !isInstalled {
		instBtn, err := gtk.ButtonNewWithLabel("Install")
		checkErr(err)
		instBtn.Show()
		buttonBox.Add(instBtn)
	} else {
		rmBtn, err := gtk.ButtonNewWithLabel("Remove")
		checkErr(err)
		rmBtn.Show()
		buttonBox.Add(rmBtn)
	}
}
