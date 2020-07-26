package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/gotk3/gotk3/gdk"

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

func onDragData(window *gtk.Window, context *gdk.DragContext, x, y int, dataPtr uintptr, info, time uint) {
	data := string(gtk.GetData(dataPtr))

	data = strings.TrimSpace(data)

	if strings.HasPrefix(data, "file://") {
		data = strings.ReplaceAll(data, "file://", "")
	} else {
		showError("unable to handle invalid file type or method of drop")
		return
	}

	if strings.HasSuffix(data, ".flatpakref") {
		log.Println("dropped flatpak refrence file -", data)
		go handleFlatpak(data)

	} else if strings.HasSuffix(data, "/recipie") {
		log.Println("dropped recipie file -", data)
		go handleRecipie(data)

	} else if strings.HasSuffix(data, "x86_64.tar.xz") {
		log.Println("dropped package file -", data)
		go handlePkg(data)
	} else {
		showError("Invalid file dropped " + data)
	}
}

func onSettingClick() {
	log.Println("clicking")
	stackPage.SetVisibleChildName("settingPage")

	getWidget("sourceEntry").(*gtk.Entry).SetText(conf.SourceURL)
	getWidget("binaryEntry").(*gtk.Entry).SetText(conf.BinaryURL)

	getWidget("recipieEntry").(*gtk.Entry).SetText(conf.RecipieDir)
	getWidget("bookEntry").(*gtk.Entry).SetText(strings.Join(conf.Repos, " "))

}

func onApplyBtnClick() {
	conf.SourceURL, _ = getWidget("sourceEntry").(*gtk.Entry).GetText()
	conf.BinaryURL, _ = getWidget("binaryEntry").(*gtk.Entry).GetText()

	conf.RecipieDir, _ = getWidget("recipieEntry").(*gtk.Entry).GetText()
	r, _ := getWidget("bookEntry").(*gtk.Entry).GetText()
	conf.Repos = strings.Split(r, " ")

}

func saveConfig() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		log.Println(err)
	}

	fptr, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln(err)
	}
	defer fptr.Close()

	_, err = fptr.Write([]byte(buf.String()))
	showError("error while saving file " + err.Error())

	fptr.Close()
}
