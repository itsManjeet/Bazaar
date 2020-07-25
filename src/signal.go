package main

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
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

func onInstallClick(widget *gtk.Button, app appData) {
	progressbar.Show()
	go installApp(widget, app)
}

func onUninstallClick(widget *gtk.Button, app appData) {

	go uninstallApp(widget, app)
}

func onUpdateClick(widget *gtk.Button, app appData) {
	go updateApp(widget, app)
}

func onAppSelect(iview *gtk.IconView, tpth *gtk.TreePath) {
	iter, err := listmodel.GetIter(tpth)
	checkErr(err)
	nameValue, _ := listmodel.GetValue(iter, 1)
	repoValue, _ := listmodel.GetValue(iter, 2)

	appName, _ := nameValue.GetString()
	repo, _ := repoValue.GetString()

	app := getapp(appName, repo)
	setupAppPage(app)

	stackPage.SetVisibleChildName("appPage")
}

func listCategory(cat string) []appData {
	catapplist := make([]appData, 0)
	for _, a := range listapps() {
		for _, c := range a.category {
			if c == cat {
				catapplist = append(catapplist, a)
				break
			}
		}
	}
	return catapplist
}

func onCategorySelect(cbox *gtk.ListBox, selrow *gtk.ListBoxRow) {

	selboxwid, _ := selrow.GetChild()

	selbox := gtk.Box{
		Container: gtk.Container{Widget: *selboxwid},
	}

	glist := selbox.GetChildren()
	lbldata := glist.NthData(1)
	lblwid := lbldata.(*gtk.Widget)
	lbl := gtk.Label{
		Widget: *lblwid,
	}

	cat, _ := lbl.GetText()

	go loadCategory(cat)
}

func onSearchChanged(searchBox *gtk.SearchEntry) {
	curapp, _ := searchBox.GetText()
	searchapp := make([]appData, 0)
	for _, a := range applist {
		if strings.Contains(a.name, curapp) {
			searchapp = append(searchapp, a)
		}
	}

	loadApps(searchapp)
}

func onRefresh(refbtn *gtk.Button) {
	refbtn.SetSensitive(false)
	go refreshData(refbtn)
}
