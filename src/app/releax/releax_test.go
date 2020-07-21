package releax

import (
	"fmt"
	"testing"
)

func TestReleaxStore(t *testing.T) {
	releaxStore := Store{
		DataDir:     "/var/lib/app/",
		RecipiePath: "/usr/recipies/",
		BinaryURL:   "https://manjeet.cloudtb.online/apps/",
		SourceURL:   "https:///github.com/itsmanjeet/packages/",
		Repositories: []string{
			"core",
			"base",
			"xorg",
			"extra",
		},
	}

	applist, err := releaxStore.ListApps()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(applist)
}
