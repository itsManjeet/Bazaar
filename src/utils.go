package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/itsmanjeet/bazaar/src/app/store"

	"github.com/gotk3/gotk3/gdk"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

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

func checkErr(e error) {
	if e != nil {
		d := dialogBox("err[software]: unexpected software side error - " + e.Error())
		d.Run()
		d.Destroy()
		log.Panic(e.Error())
	}
}

func getWidget(widgetName string) glib.IObject {
	obj, err := builder.GetObject(widgetName)
	checkErr(err)

	return obj
}

func dialogBox(msg string) *gtk.MessageDialog {
	dialog := gtk.MessageDialogNew(
		nil,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK_CANCEL,
		msg,
	)

	return dialog
}

func categoryLabel(name string) *gtk.Label {
	lbl, err := gtk.LabelNew(name)
	checkErr(err)

	lbl.SetMarginBottom(16)
	lbl.SetMarginTop(16)
	lbl.SetMarginStart(45)
	lbl.SetMarginEnd(45)
	lbl.Show()
	return lbl
}

func pixbuftype() glib.Type {
	pixbuf, err := gdk.PixbufNew(gdk.COLORSPACE_RGB, true, 8, 64, 64)
	checkErr(err)
	return pixbuf.TypeFromInstance()
}

func appendApp(app store.App) {
	iter := listmodel.Append()
	icon := app.Icon(64)

	listmodel.SetValue(iter, 0, icon)
	listmodel.SetValue(iter, 1, app.Name())
}

func loadApps(apps []store.App) {
	glib.IdleAdd(listmodel.Clear)
	for _, a := range apps {
		glib.IdleAdd(appendApp, a)
	}
}

func execApp(btn *gtk.Button, cmd string) {
	exec.Command("pkexec", "--disable-internal-agent", "exo-open", cmd).Start()
}
