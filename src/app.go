package main

import (
	"github.com/itsmanjeet/bazaar/src/app/flatpak"
	"github.com/itsmanjeet/bazaar/src/app/releax"
	"github.com/itsmanjeet/bazaar/src/app/store"
)

func listApps() []store.App {
	stores = make([]store.Store, 0)
	releaxStore = releax.Store{
		DataDir:      "/var/lib/app",
		Repositories: []string{"core", "base", "xorg", "extra"},
		RecipiePath:  "/usr/recipies/",
		SourceURL:    "https://github.com/itsmanjeet/packages/",
		BinaryURL:    "https://manjeet.cloudtb.online/apps/",
	}

	flatpakStore := flatpak.Store{
		AppStream: "/var/lib/flatpak/appstream/flathub/x86_64/active/appstream.xml",
	}

	stores = append(stores, releaxStore)
	stores = append(stores, flatpakStore)

	applist = make([]store.App, 0)

	for _, s := range stores {
		apl, err := s.ListApps()
		if err != nil {
			checkErr(err)
			continue
		}

		applist = append(applist, apl...)
	}

	return applist
}
