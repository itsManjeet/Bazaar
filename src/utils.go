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

func categoryLabel(name string, icon string) *gtk.Box {

	gbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)

	lbl, err := gtk.LabelNew(name)
	checkErr(err)

	gbox.SetMarginBottom(16)
	gbox.SetMarginTop(16)
	gbox.SetMarginStart(45)
	gbox.SetMarginEnd(45)

	ic, _ := gtk.ImageNewFromIconName(icon, 0)
	ic.Show()
	lbl.Show()
	gbox.PackStart(ic, false, false, 1)
	gbox.PackStart(lbl, true, true, 1)

	gbox.ShowAll()
	return gbox
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
	glib.IdleAdd(listmodel.Clear)
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
	backbtn := getWidget("backButton").(*gtk.Button)

	glib.IdleAdd(backbtn.SetSensitive, true)

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

		fptr, err := ioutil.ReadFile(conf.DataDir + "/" + app.name + "/timestamp")
		if err != nil {
			instLabel.SetText(err.Error())
		} else {
			instLabel.SetText(string(fptr))
		}

		apflst, err := ioutil.ReadFile(conf.DataDir + "/" + app.name + "/files")
		if err == nil {
			lsbuf.SetText(string(apflst))
		}

	}

	btn.Show()
	buttonBox.Add(btn)
}

func loadCategory(cat string) {
	if val, ok := cacheData[cat]; ok {
		log.Println("loading for cache")
		glib.IdleAdd(loadApps, val)

		return
	}

	acl := make([]appData, 0)
	if cat == "Market" {
		for _, a := range listapps() {
			if a.repo == "extra" {
				acl = append(acl, a)
			}
		}
	} else if cat == "Must Have" {
		for _, a := range []string{
			"Accessories", "Graphics", "Internet", "Multimedia", "Office",
		} {
			acl = append(acl, listCategory(a)...)
		}
	} else if cat == "Personalize" {

		for _, a := range []string{
			"Customizations", "Plugins",
		} {
			acl = append(acl, listCategory(a)...)
		}
	} else if cat == "Developer" {
		for _, a := range []string{
			"Development", "Libraries", "Library",
		} {
			acl = append(acl, listCategory(a)...)
		}
	} else if cat == "System" {
		acl = listCategory(cat)
	} else if cat == "Games" {
		acl = listCategory(cat)
	}

	log.Println("Cat :", cat)

	cacheData[cat] = make([]appData, 0)
	cacheData[cat] = append(cacheData[cat], acl...)
	glib.IdleAdd(loadApps, acl)

}
