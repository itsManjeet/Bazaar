package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func uninstallApp(widget *gtk.Button, app appData) bool {
	log.Println("Uninstalling", app.name)

	glib.IdleAdd(widget.SetSensitive, false)
	backbtn := getWidget("backButton").(*gtk.Button)
	glib.IdleAdd(backbtn.SetSensitive, false)
	glib.IdleAdd(progressbar.Show)
	glib.IdleAdd(progressbar.SetText, "Uninstalling")
	togo := true
	go func() {
		for togo {
			glib.IdleAdd(progressbar.Pulse)
			time.Sleep(time.Millisecond * 50)
		}
	}()
	out, err := exec.Command("sys-app", "rm", app.name, "--no-ask").Output()
	togo = false
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
