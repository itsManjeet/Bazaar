package releax

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// Recipie app data
type Recipie struct {
	name        string
	version     string
	release     string
	description string
	url         string
	license     string
	apptype     string
	depends     []string
	icon        string
	filepath    string
	categories  []string
}

// Name return name
func (r *Recipie) Name() string {
	return r.name
}

func (r *Recipie) ID() string {
	return r.name
}

func (r *Recipie) Version() string {
	return r.version
}

func (r *Recipie) Release() string {
	return r.release
}

func (r *Recipie) Description() string {
	return r.description
}

func (r *Recipie) Depends() []string {
	return r.depends
}

func (r *Recipie) Categories() []string {
	return r.categories
}

func (r *Recipie) Url() string {
	return r.url
}

func (r *Recipie) License() string {
	return r.license
}

func (r *Recipie) Store() string {
	return "releax"
}

func (r *Recipie) Icon(size uint) *gdk.Pixbuf {
	ict, _ := gtk.IconThemeGetDefault()
	var err error
	var pixbuf *gdk.Pixbuf
	if len(r.icon) != 0 {
		pixbuf, err = ict.LoadIcon(r.icon, int(size), 0)
		if err == nil {
			return pixbuf
		}
	}

	if exists(r.filepath + "/icon") {
		pixbuf, err = gdk.PixbufNewFromFileAtSize(r.filepath+"/icon", int(size), int(size))
		if err == nil {
			return pixbuf
		}
	}

	pixbuf, err = ict.LoadIcon(r.name, int(size), 0)
	if err == nil {
		return pixbuf
	}

	if strings.HasPrefix(r.name, "lib") {
		pixbuf, err = ict.LoadIcon("application-x-sharedlib", int(size), 0)
		if err == nil {
			return pixbuf
		}
	} else {
		pixbuf, err = ict.LoadIcon("application-x-pak", int(size), 0)
		if err == nil {
			return pixbuf
		}
	}
	return nil
}

func (r *Recipie) PackageFile() string {
	return fmt.Sprintf("%s-%s-%s-x86_64.tar.xz", r.name, r.version, r.release)
}

// NewFrom return new recipie from paring recipie file
func NewFrom(path string) (*Recipie, error) {

	rcpFile, err := os.OpenFile(path+"/recipie", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer rcpFile.Close()

	var rcp Recipie
	rcp.filepath = path

	scanner := bufio.NewScanner(rcpFile)
	for scanner.Scan() {
		curline := scanner.Text()
		curline = strings.TrimSpace(curline)
		if len(curline) == 0 {
			continue
		}

		if curline[0] == '#' {
			if strings.HasPrefix(curline, "# Description") {
				rcp.description = getCommentedValue(curline, "Description")
			} else if strings.HasPrefix(curline, "# URL") {
				rcp.url = getCommentedValue(curline, "URL")
			} else if strings.HasPrefix(curline, "# License") {
				rcp.license = getCommentedValue(curline, "License")
			} else if strings.HasPrefix(curline, "# Depends on") {
				rcp.depends = strings.Split(getCommentedValue(curline, "Depends on"), " ")
			} else if strings.HasPrefix(curline, "# Icon") {
				rcp.icon = getCommentedValue(curline, "Icon")
			} else if strings.HasPrefix(curline, "# Category") {
				rcp.categories = strings.Split(getCommentedValue(curline, "Category"), " ")
			}
		} else if strings.Contains(curline, "=") {
			if strings.Contains(curline, "version=") {
				rcp.version = getAssignValue(curline, "version")
			} else if strings.Contains(curline, "release=") {
				rcp.release = getAssignValue(curline, "release")
			} else if strings.Contains(curline, "name=") {
				rcp.name = getAssignValue(curline, "name")
			} else if strings.Contains(curline, "type=") {
				rcp.apptype = getAssignValue(curline, "type")
			}
		}
	}

	return &rcp, nil
}

func getCommentedValue(curline, cmt string) string {
	if strings.HasPrefix(curline, "# "+cmt) {
		return strings.TrimSpace(curline[strings.Index(curline, ":")+1 : len(curline)])
	}
	return string("")
}

func getAssignValue(curline, variable string) string {
	if strings.Contains(curline, "=") {
		rindex := strings.Index(curline, "=")
		vari := curline[0:rindex]
		if vari == variable {
			return strings.TrimSpace(curline[rindex+1 : len(curline)])
		}

	}
	return string("")
}
