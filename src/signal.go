package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/itsmanjeet/bazaar/src/app/store"
)

func onBackClick() {
	clearAppPage()
	stackPage.SetVisibleChildName("marketPage")
}

func clearAppPage() {
	logoImage := getWidget("logoImage").(*gtk.Image)
	buttonBox := getWidget("buttonBox").(*gtk.Box)
	logoImage.Clear()

	progressbar.SetText("")
	progressbar.SetFraction(0.0)
	progressbar.Hide()

	buttonBox.GetChildren().Foreach(func(widget interface{}) {
		widget.(*gtk.Widget).Destroy()
	})
}

func onInstallClick(widget *gtk.Button, app store.App) {
	progressbar.Show()
}

func onUninstallClick(widget *gtk.Button, app store.App) {

}

func onUpdateClick(widget *gtk.Button, app store.App) {
}

func onAppSelect(iview *gtk.IconView, tpth *gtk.TreePath) {
	iter, err := listmodel.GetIter(tpth)
	checkErr(err)
	nameValue, _ := listmodel.GetValue(iter, 1)
	appName, _ := nameValue.GetString()

	var app store.App

	for _, s := range stores {
		app, err = s.GetApp(appName)
		if err == nil {
			break
		}
	}

	if app == nil {
		checkErr(errors.New("Unknown error: unable to find listed app"))
	}
	setupPage(app)
	stackPage.SetVisibleChildName("appPage")
}

func onCategorySelect(cbox *gtk.ListBox, selrow *gtk.ListBoxRow) {
	catlbl, _ := selrow.GetChild()
	lbl, _ := gtk.WidgetToLabel(catlbl)
	lable := lbl.(*gtk.Label)

	cat, _ := lable.GetText()
	if cat == "All" {
		glib.IdleAdd(loadApps, applist)
	} else {
		catapplist := make([]store.App, 0)
		for _, a := range applist {
			for _, c := range a.Categories() {
				if strings.ToLower(c) == strings.ToLower(cat) {
					catapplist = append(catapplist, a)
					break
				}
			}
		}

		glib.IdleAdd(loadApps, catapplist)
	}

}

func onSearchChanged(searchBox *gtk.SearchEntry) {
	curapp, _ := searchBox.GetText()
	searchapp := make([]store.App, 0)
	for _, a := range applist {
		if strings.Contains(strings.ToLower(a.Name()), strings.ToLower(curapp)) {
			searchapp = append(searchapp, a)
		}
	}

	loadApps(searchapp)
}

func onRefresh(refbtn *gtk.Button) {
	refbtn.SetSensitive(false)
	//go refreshData(refbtn)
}

func onDragDrop(window *gtk.Window, ctx *gdk.DragContext, x, y int, data uintptr, m int, t uint) {
	fmt.Println(x, y)
	fmt.Println("Reciving data")
}
