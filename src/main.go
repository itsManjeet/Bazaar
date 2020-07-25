package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	checkErr(err)

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		showError(err.Error())
	}

	application.Connect("startup", func() {
		log.Println("application starting")
	})

	application.Connect("activate", func() {
		log.Println("application activating")

		builder, err = gtk.BuilderNewFromFile(uiFile)
		checkErr(err)

		signals := map[string]interface{}{
			"onAppSelect":      onAppSelect,
			"onBackClick":      onBackClick,
			"onCategorySelect": onCategorySelect,
			"onSearchChanged":  onSearchChanged,
			"onRefresh":        onRefresh,
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

		appIconView.SetPixbufColumn(0)
		appIconView.SetTextColumn(1)

		buildCat := func() {
			glib.IdleAdd(catListBox.Add, categoryLabel("Market", "amarok_cart_view"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Must Have", "folder-bookmark"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Personalize", "draw-brush"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Games", "folder-games"))
			glib.IdleAdd(catListBox.Add, categoryLabel("Developer", "format-text-code"))
			glib.IdleAdd(catListBox.Add, categoryLabel("System", "configure"))
		}

		go buildCat()

		icontheme, err = gtk.IconThemeGetDefault()
		checkErr(err)
		applist = listapps()
		window.Show()
		application.AddWindow(window)

	})

	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	os.Exit(application.Run(os.Args))
}
