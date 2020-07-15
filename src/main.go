package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	checkErr(err)

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
			"onDragDrop":       onDragDrop,
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

		for _, a := range categories {
			catListBox.Add(categoryLabel(a))
		}

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
