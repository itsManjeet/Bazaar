package main

import (
	"log"
	"os/exec"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func uninstallApp(widget *gtk.Button, app appData) bool {
	log.Println("Uninstalling", app.name)

	glib.IdleAdd(widget.SetSensitive, false)
	glib.IdleAdd(progressbar.Show)
	glib.IdleAdd(progressbar.SetText, "Uninstalling")
	out, err := exec.Command("sys-app", "rm", app.name).Output()

	if err != nil {
		glib.IdleAdd(showError, string(out)+"\nError: "+err.Error())
	} else {
		glib.IdleAdd(progressbar.SetFraction, 1.0)
		glib.IdleAdd(progressbar.SetText, "removed successfully")
	}

	glib.IdleAdd(clearAppPage)
	glib.IdleAdd(setupAppPage, app)
	return true
}
