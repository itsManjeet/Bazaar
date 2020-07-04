package main

import (
	"io/ioutil"
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

func onCategorySelect(cbox *gtk.ListBox, selrow *gtk.ListBoxRow) {
	catlbl, _ := selrow.GetChild()
	lbl, _ := gtk.WidgetToLabel(catlbl)
	lable := lbl.(*gtk.Label)

	cat, _ := lable.GetText()

	if cat == "All" {
		loadApps(listapps())
	} else if cat == "Installed" {
		instdir, err := ioutil.ReadDir(datadir)
		checkErr(err)
		aplst := make([]appData, 0)
		for _, z := range instdir {
			if !z.IsDir() {
				continue
			}
			a, err := getFromAppList(z.Name())
			if err != nil {
				continue
			}
			aplst = append(aplst, a)
		}
		loadApps(aplst)
	} else {
		catapplist := make([]appData, 0)
		for _, a := range listapps() {
			for _, c := range a.category {
				if c == cat {
					catapplist = append(catapplist, a)
					break
				}
			}
		}

		loadApps(catapplist)
	}

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
