package flatpak

import (
	"fmt"
	"testing"
)

func TestFlatpak(t *testing.T) {
	flatpakStore := Store{
		AppStream: "/var/lib/flatpak/appstream/flathub/x86_64/active/appstream.xml",
	}

	applist, err := flatpakStore.ListApps()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(applist)
}
