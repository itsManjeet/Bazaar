package main

import (
	"os/exec"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func refreshData(refbtn *gtk.Button) {
	cmd := exec.Command("sys-app", "rf")
	s := true
	glib.IdleAdd(refProgress.SetVisible, true)
	go func() {
		for s {
			glib.IdleAdd(refProgress.Pulse, "refreshing database")
			time.Sleep(time.Second)
		}

	}()
	if err := cmd.Start(); err != nil {
		glib.IdleAdd(showError, err.Error())
		glib.IdleAdd(refbtn.SetSensitive, true)
		glib.IdleAdd(refProgress.SetVisible, false)
		s = false
		return
	}

	if err := cmd.Wait(); err != nil {
		glib.IdleAdd(showError, err.Error())
		glib.IdleAdd(refbtn.SetSensitive, true)
		glib.IdleAdd(refProgress.SetVisible, false)
		s = false
		return
	}

	cmdout, _ := exec.Command("sys-app", "dry-check-update").Output()
	if len(cmdout) != 0 {
		a := strings.Split(string(cmdout), "\n")
		upda := ""
		for _, x := range a {
			if len(x) == 0 {
				continue
			}
			if x[0] == '\033' {
				continue
			}
			upda += " " + x
		}

		applist = listapps()

		apd := make([]appData, 0)
		for _, x := range strings.Split(upda, " ") {
			ap, err := getFromAppList(x)
			if err != nil {
				continue
			}
			if ap.name == "double-conversion" {
				continue
			}
			apd = append(apd, ap)
		}
	}

	glib.IdleAdd(refProgress.SetVisible, false)
	glib.IdleAdd(refbtn.SetSensitive, true)
	s = false
	return

}
