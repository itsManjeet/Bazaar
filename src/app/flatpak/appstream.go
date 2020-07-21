package flatpak

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// AppStream of flatpak
type AppStreamData struct {
	XMLName    string      `xml:"components"`
	Components []Component `xml:"component"`
}

// Component of flatpak appstream
type Component struct {
	ID         string `xml:"id"`
	Xmlname    string `xml:"name"`
	Xmlsummary string `xml:"summary"`
	Xmllicense string `xml:"project_license"`
}

func (c Component) Name() string {
	return c.Xmlname
}

func (c Component) Version() string {
	return ""
}

func (c Component) Release() string {
	return ""
}

func (c Component) Description() string {
	return c.Xmlsummary
}

func (c Component) PackageFile() string {
	return ""
}

func (c Component) Depends() []string {
	return nil
}

func (c Component) Icon(size uint) *gdk.Pixbuf {
	appstreamAddr := "/var/lib/flatpak/appstream/flathub/x86_64/active/icons/"
	addr := fmt.Sprintf("%s/%dx%d/%s.png", appstreamAddr, int(size), int(size), c.ID)
	pixbuf, err := gdk.PixbufNewFromFileAtSize(addr, int(size), int(size))
	if err != nil {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(addr, 64, 64)
		if err == nil {
			return pixbuf
		}
		//log.Println(c.Name()+"->Icon() failed: ", err)
		ict, _ := gtk.IconThemeGetDefault()
		pixbuf, err = ict.LoadIcon(strings.ToLower(c.Name()), int(size), 0)
		if err == nil {
			return pixbuf
		}

		pb, _ := ict.LoadIcon("application-x-pak", int(size), 0)
		return pb
	}
	return pixbuf
}
