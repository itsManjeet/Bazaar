package main

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/gotk3/gotk3/glib"
	"github.com/sqp/vte/vte.gtk3"
)

var term *vte.Terminal

func onBackFromProcess() {
	term.Destroy()
	glib.IdleAdd(stackPage.SetVisibleChildName, "marketPage")
}

func onDone() {
	backBtn := getWidget("backBtn").(*gtk.Button)
	glib.IdleAdd(backBtn.SetSensitive, true)
}

func doProcess(cmd []string, dir string) error {
	glib.IdleAdd(stackPage.SetVisibleChildName, "processPage")

	backBtn := getWidget("backBtn").(*gtk.Button)
	glib.IdleAdd(backBtn.SetSensitive, false)

	term = vte.NewTerminal()

	processBox := getWidget("processBox").(*gtk.Box)

	glib.IdleAdd(processBox.Add, term)

	term.SetMarginTop(50)
	term.Show()
	term.SetVExpand(true)
	term.GrabFocus()

	term.Connect("child-exited", onDone)

	_, err := term.ExecSync(dir, cmd, nil)
	if err != nil {
		glib.IdleAdd(showError, err.Error())
		return err
	}

	return nil
}
