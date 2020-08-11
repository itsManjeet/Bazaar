package main

import (
	"github.com/gotk3/gotk3/gtk"
)

const (
	appID = "in.releax.bazaar"
)

// Config bazaar configuration
type Config struct {
	SourceURL string
	BinaryURL string

	DataDir    string
	RecipieDir string
	Repos      []string
}

var conf Config
var cacheData map[string][]appData
var (
	builder     *gtk.Builder
	icontheme   *gtk.IconTheme
	listmodel   *gtk.ListStore
	stackPage   *gtk.Stack
	progressbar *gtk.ProgressBar
	refProgress *gtk.ProgressBar

	applist []appData
)
