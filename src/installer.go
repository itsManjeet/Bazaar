package main

import (
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gtk"
)

func showError(err string) {
	d := dialogBox(err)
	d.Run()
	d.Destroy()
}

func installApp(widget *gtk.Button, app appData) bool {
	log.Println("Installing", app.name)

	dwndlr := downloader{
		prgsfunc: func(prgs int) {
			glib.IdleAdd(progressbar.SetText, "Downloading "+strconv.Itoa(prgs)+"%...")
			parprgs := (float32(prgs) / 100) * 0.6
			glib.IdleAdd(progressbar.SetFraction, parprgs)
		},
	}
	backbtn := getWidget("backButton").(*gtk.Button)
	glib.IdleAdd(backbtn.SetSensitive, false)
	glib.IdleAdd(widget.SetSensitive, false)

	glib.IdleAdd(progressbar.Show)

	installer := func(a appData) bool {
		appname := a.name + "-" + a.version + "-" + a.release + "-x86_64.tar.xz"
		apppath := "/var/cache/build/" + appname
		glib.IdleAdd(progressbar.SetFraction, 0.0)
		dwndlr.current = 0
		dwndlr.total = 0
		if !exists(apppath) {
			if err := dwndlr.download(apppath, conf.BinaryURL+"/"+appname); err != nil {
				glib.IdleAdd(progressbar.SetText, "Error while downoading "+err.Error())
				time.Sleep(time.Second * 1)
				return false
			}
		}

		glib.IdleAdd(progressbar.SetText, "Installing")
		togo := true
		go func() {
			for togo {
				glib.IdleAdd(progressbar.Pulse)
				time.Sleep(time.Millisecond * 50)
			}
		}()
		out, err := exec.Command("sys-app", "in", a.name, "--no-ask").Output()
		togo = false
		if err != nil {
			glib.IdleAdd(showError, string(out)+"\nError: "+err.Error())
			return false

		}
		glib.IdleAdd(progressbar.SetFraction, 1.0)
		glib.IdleAdd(progressbar.SetText, "Installed Success")
		return true
	}

	deps := app.getDepends()
	if len(deps) != 0 {
		for _, a := range deps {
			if !installer(a) {
				glib.IdleAdd(clearAppPage)
				glib.IdleAdd(setupAppPage, app)
				return false
			}
		}
	}

	status := installer(app)

	glib.IdleAdd(clearAppPage)
	glib.IdleAdd(setupAppPage, app)
	return status
}
