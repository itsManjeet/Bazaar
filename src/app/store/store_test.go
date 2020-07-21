package store_test

import (
	"fmt"
	"testing"

	"github.com/itsmanjeet/app/src/releax"

	"github.com/itsmanjeet/app/src/store"
)

func TestStore(t *testing.T) {
	stores := make([]store.Store, 0)
	releaxStore := releax.Store{
		DataDir:     "/var/lib/app",
		RecipiePath: "/usr/recipies",
		Repositories: []string{
			"core", "extra", "base", "xorg",
		},
	}

	stores = append(stores, releaxStore)

	apps, err := stores[0].ListApps()
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range apps {
		fmt.Println(a.Name())
	}
}
