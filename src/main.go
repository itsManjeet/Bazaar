package main

import (
	"log"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"gopkg.in/ini.v1"
)

func main() {

	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	checkErr(err)

	cfg, err := ini.Load(configFile)
	checkErr(err)

	conf.SourceURL = cfg.Section("url").Key("source").String()
	conf.BinaryURL = cfg.Section("url").Key("binary").String()

	conf.DataDir = cfg.Section("dir").Key("data").String()
	conf.RecipieDir = cfg.Section("dir").Key("recipies").String()
	conf.Repos = strings.Split(cfg.Section("default").Key("repo").String(), " ")

	cacheData = make(map[string][]appData)
	application.Connect("startup", func() {
		log.Println("application starting")
	})

	application.Connect("activate", func() {
		log.Println("application activating")

		builder, err = gtk.BuilderNewFromFile(uiFile)
		checkErr(err)

		signals := map[string]interface{}{
			"onAppSelect":       onAppSelect,
			"onBackClick":       onBackClick,
			"onCategorySelect":  onCategorySelect,
			"onSearchChanged":   onSearchChanged,
			"onRefresh":         onRefresh,
			"onBackFromProcess": onBackFromProcess,
			"onSettingClick":    onSettingClick,
			"onApplyBtnClick":   onApplyBtnClick,
		}

		builder.ConnectSignals(signals)

		window := getWidget("mainWindow").(*gtk.Window)
		catListBox := getWidget("catListBox").(*gtk.ListBox)
		appIconView := getWidget("appIconView").(*gtk.IconView)
		stackPage = getWidget("stackPage").(*gtk.Stack)
		progressbar = getWidget("progressBar").(*gtk.ProgressBar)
		refProgress = getWidget("refProgress").(*gtk.ProgressBar)

		listmodel, err = gtk.ListStoreNew(pixbuftype(), glib.TYPE_STRING, glib.TYPE_STRING)
		checkErr(err)

		appIconView.SetModel(listmodel)

		window.Connect("drag-data-received", onDragData)
		// Setup drag and drop
		var targets []gtk.TargetEntry
		te, err := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
		checkErr(err)

		targets = append(targets, *te)
		window.DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			targets,
			gdk.ACTION_COPY,
		)

		appIconView.SetPixbufColumn(0)
		appIconView.SetTextColumn(1)

		buildCat := func() {
			glib.IdleAdd(catListBox.Add, categoryLabel("Market", "amarok_cart_view"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Must Have", "folder-bookmark"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Personalize", "draw-brush"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Games", "folder-games"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Developer", "format-text-code"))
			glib.IdleAdd(catListBox.Add, categoryLabel("System", "configure"))

			glib.IdleAdd(catListBox.SelectRow, catListBox.GetRowAtY(-1))
		}
		icontheme, err = gtk.IconThemeGetDefault()
		checkErr(err)
		applist = listapps()
		buildCat()

		window.Show()
		application.AddWindow(window)

	})

	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	os.Exit(application.Run(os.Args))
}
