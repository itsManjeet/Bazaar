package main

import (
	"io/ioutil"
	"log"
	"os"

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

func appendApp(app appData) {
	iter := listmodel.Append()
	icon := app.geticon(64)

	listmodel.SetValue(iter, 0, icon)
	listmodel.SetValue(iter, 1, app.name)
	listmodel.SetValue(iter, 2, app.repo)
}

func loadApps(apps []appData) {
	listmodel.Clear()
	for _, a := range apps {
		glib.IdleAdd(appendApp, a)
	}
}

func setupAppPage(app appData) {
	logoImage := getWidget("logoImage").(*gtk.Image)
	nameLabel := getWidget("nameLabel").(*gtk.Label)
	descLabel := getWidget("descLabel").(*gtk.Label)
	buttonBox := getWidget("buttonBox").(*gtk.Box)
	instdData := getWidget("instdData").(*gtk.Grid)
	reqLabel := getWidget("reqLabel").(*gtk.Label)
	infoTab := getWidget("infoTab").(*gtk.Notebook)

	instLabel := getWidget("instLabel").(*gtk.Label)
	licenseLabel := getWidget("licenseLabel").(*gtk.Label)
	urlLabel := getWidget("urlLabel").(*gtk.Label)

	fileList := getWidget("fileList").(*gtk.TextView)

	infoTab.SetCurrentPage(0)
	deps := app.getDepends()
	if len(deps) == 0 {
		reqLabel.SetMarkup("<span foreground=\"green\"><b>✔️ ok</b></span>")
	} else {
		depstr := ""
		for _, x := range deps {
			depstr += x.name + " "
		}

		reqLabel.SetText(depstr)
	}

	logoImage.SetFromPixbuf(app.geticon(128))
	descLabel.SetText(app.description)

	licenseLabel.SetText(app.license)
	urlLabel.SetText(app.url)

	lsbuf, _ := fileList.GetBuffer()
	lsbuf.SetText("")

	var btn *gtk.Button
	if !app.isInstalled() {
		nameLabel.SetText(app.name + " " + app.version + " " + app.release)
		btn, _ = gtk.ButtonNewWithLabel("install")
		btn.Connect("clicked", onInstallClick, app)
		instdData.SetVisible(false)

	} else {

		instver := getInstVer(app)
		instrel := getInstRel(app)
		update := false
		labelstring := app.name + " "
		if instver != app.version {
			labelstring += "[" + instver + "->" + app.version + "] "
			update = true
		} else {
			labelstring += instver + " "
		}

		if instrel != app.release {
			labelstring += "[" + instrel + "->" + app.release + "]"
			update = true
		} else {
			labelstring += instrel
		}

		nameLabel.SetText(labelstring)

		if update {
			upbtn, _ := gtk.ButtonNewWithLabel("update")
			upbtn.Connect("clicked", onUpdateClick, app)
			upbtn.Show()
			buttonBox.Add(upbtn)
		}
		btn, _ = gtk.ButtonNewWithLabel("uninstall")
		btn.Connect("clicked", onUninstallClick, app)
		instdData.SetVisible(true)

		fptr, err := ioutil.ReadFile(datadir + app.name + "/timestamp")
		if err != nil {
			instLabel.SetText(err.Error())
		} else {
			instLabel.SetText(string(fptr))
		}

		apflst, err := ioutil.ReadFile(datadir + app.name + "/files")
		if err == nil {
			lsbuf.SetText(string(apflst))
		}

	}

	btn.Show()
	buttonBox.Add(btn)
}
