package main

import (
	"github.com/gotk3/gotk3/gtk"
)

const (
	appID   = "in.releax.bazaar"
	repodir = "/usr/recipies/" // Repository, example "/usr/rcipies/core/acl/recipie"
	datadir = "/var/lib/app/"  // Database, contain information of installed apps
	repourl = "https://manjeet.cloudtb.online/apps/"
)

var categories = []string{
	"All",
	"Accessories", "Development", "Graphics", "Internet", "Games", "Multimedia",
	"Office", "Customizations", "Plugins", "System", "Libraries", "Installed",
}

var (
	builder     *gtk.Builder
	icontheme   *gtk.IconTheme
	listmodel   *gtk.ListStore
	stackPage   *gtk.Stack
	progressbar *gtk.ProgressBar
	refProgress *gtk.ProgressBar

	applist []appData
)
