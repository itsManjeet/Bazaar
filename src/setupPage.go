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
	appDesc := getWidget("descLabel").(*gtk.Label)

	imagePix := app.Icon(128)
	if imagePix != nil {
		appImage.SetFromPixbuf(imagePix)
	}
	appName.SetText(fmt.Sprintf("%s %s %s", app.Name(), app.Version(), app.Release()))
	appDesc.SetText(app.Description())
}
