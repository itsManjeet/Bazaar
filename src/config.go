package main

import (
	"github.com/gotk3/gotk3/gtk"
)

const (
	appID = "in.releax.bazaar"
)

// Config bazaar configuration
type Config struct {
	SourceURL string `toml:"source"`
	BinaryURL string `toml:"binary"`

	DataDir    string   `toml:"data"`
	RecipieDir string   `toml:"recipies"`
	Repos      []string `toml:"repo"`
}

var conf Config

var (
	builder     *gtk.Builder
	icontheme   *gtk.IconTheme
	listmodel   *gtk.ListStore
	stackPage   *gtk.Stack
	progressbar *gtk.ProgressBar
	refProgress *gtk.ProgressBar

	applist []appData
)
