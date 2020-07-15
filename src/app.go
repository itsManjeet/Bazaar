package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gdk"
)

type appData struct {
	name        string
	repo        string
	version     string
	release     string
	description string
	url         string
	icon        string
	license     string
	category    []string
	depends     []string
}

type desktopFile struct {
	name        string
	command     string
	icon        string
	desktopfile string
}

func parseDesktopFile(add string) (desktopFile, error) {
	file, err := os.Open(add)
	checkErr(err)
	defer file.Close()

	var desk desktopFile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "=") {
			data := strings.Split(scanner.Text(), "=")
			if data[0] == "Name" {
				desk.name = data[1]
			} else if data[0] == "Exec" {
				desk.command = data[1]
			} else if data[0] == "Icon" {
				desk.icon = data[1]
			}
		}
	}

	desk.desktopfile = add

	return desk, nil
}

func (a *appData) getDepends() []appData {

	apl := make([]appData, 0)

	for _, p := range a.depends {
		app, err := getFromAppList(p)
		if err != nil {
			continue
		}
		if app.isInstalled() || aplcontain(apl, app.name) {
			continue
		} else {
			for _, a := range app.getDepends() {
				if app.isInstalled() || aplcontain(apl, a.name) {
					continue
				}
				apl = append(apl, a)
			}
			apl = append(apl, app)
		}
	}
	return apl
}

func aplcontain(apl []appData, appname string) bool {
	for _, a := range apl {
		if a.name == appname {
			return true
		}
	}
	return false
}

func getFromAppList(appname string) (appData, error) {
	for _, a := range applist {
		if a.name == appname {
			return getapp(a.name, a.repo), nil
		}
	}
	return appData{}, errors.New("Not found")
}

func getDesktopFile(app appData) ([]desktopFile, error) {
	if !app.isInstalled() {
		return []desktopFile{}, errors.New("app is not installed")
	}

	desktopFiles := make([]desktopFile, 0)

	file, err := os.Open(path.Join(datadir, app.name, "files"))
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasSuffix(scanner.Text(), ".desktop") {
			desk, err := parseDesktopFile("/" + scanner.Text())
			if err == nil {
				desktopFiles = append(desktopFiles, desk)
			}

		}
	}

	return desktopFiles, nil
}

func getapp(name string, repo string) appData {
	app := appData{
		name: name,
		repo: repo,
	}

	file, err := os.Open(path.Join(repodir, repo, name, "recipie"))
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curline := scanner.Text()
		curline = strings.TrimSpace(curline)
		if len(curline) == 0 {
			continue
		}
		if curline[0] == '#' {
			if strings.HasPrefix(curline, "# Description") {
				app.description = getcmntval(curline, "Description")
			} else if strings.HasPrefix(curline, "# URL") {
				app.url = getcmntval(curline, "URL")
			} else if strings.HasPrefix(curline, "# License") {
				app.license = getcmntval(curline, "License")
			} else if strings.HasPrefix(curline, "# Depends on") {
				app.depends = strings.Split(getcmntval(curline, "Depends on"), " ")
			} else if strings.HasPrefix(curline, "# Icon") {
				app.icon = getcmntval(curline, "# Icon")
			}
		} else if strings.Contains(curline, "=") {
			if strings.Contains(curline, "version=") {
				app.version = getasnval(curline, "version")
			}
			if strings.Contains(curline, "release=") {
				app.release = getasnval(curline, "release")
			}
		}
	}

	return app
}

func getcmntval(curline, cmnt string) string {
	if strings.HasPrefix(curline, "# "+cmnt) {
		return strings.TrimSpace(curline[strings.Index(curline, ":")+1 : len(curline)])
	}
	return string("")
}

func getasnval(curline, variable string) string {
	if strings.Contains(curline, "=") {
		rindex := strings.Index(curline, "=")
		vari := curline[0:rindex]
		if vari == variable {
			return strings.TrimSpace(curline[rindex+1 : len(curline)])
		}

	}
	return string("")
}

func (a appData) geticon(size int) *gdk.Pixbuf {
	var pixbuf *gdk.Pixbuf
	var err error
	recpath := path.Join(repodir, a.repo, a.name)
	if exists(recpath + "/icon") {
		pixbuf, err = gdk.PixbufNewFromFileAtSize(recpath+"/icon", size, size)
	} else if len(a.icon) != 0 {
		pixbuf, err = icontheme.LoadIcon(a.icon, size, 0)
	} else {
		pixbuf, err = icontheme.LoadIcon(a.name, size, 0)
	}
	if err != nil {
		pixbuf, err = icontheme.LoadIcon("application-x-pak", size, 0)
		if err != nil {
			return nil
		}
	}

	return pixbuf
}

func (a appData) isInstalled() bool {
	return exists(path.Join(datadir, a.name))
}

func getInstVer(app appData) string {
	if !app.isInstalled() {
		return string("")
	}

	file, err := os.Open(path.Join(datadir, app.name, "info"))
	checkErr(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curline := scanner.Text()
		if strings.Contains(curline, ":") {
			rind := strings.Index(curline, ":")
			vari := strings.TrimSpace(curline[0:rind])
			if vari == "version" {
				return strings.TrimSpace(curline[rind+1 : len(curline)])
			}
		}
	}
	return string("")
}

func getInstRel(app appData) string {
	if !app.isInstalled() {
		return string("")
	}

	file, err := os.Open(path.Join(datadir, app.name, "info"))
	checkErr(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curline := scanner.Text()
		if strings.Contains(curline, ":") {
			rind := strings.Index(curline, ":")
			vari := strings.TrimSpace(curline[0:rind])
			if vari == "release" {
				return strings.TrimSpace(curline[rind+1 : len(curline)])
			}
		}
	}
	return string("")
}

func listapps() []appData {

	applist := make([]appData, 0)

	repo, err := ioutil.ReadDir(repodir)
	checkErr(err)

	for _, r := range repo {
		if !r.IsDir() {
			continue
		}

		appdir, err := ioutil.ReadDir(path.Join(repodir, r.Name()))
		checkErr(err)
		for _, a := range appdir {
			recpp := path.Join(repodir, r.Name(), a.Name(), "recipie")
			if !a.IsDir() || !exists(recpp) {
				continue
			}

			app := appData{
				name: a.Name(),
				repo: r.Name(),
			}
			file, err := os.Open(recpp)
			checkErr(err)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				curline := scanner.Text()
				curline = strings.TrimSpace(curline)
				if len(curline) == 0 {
					continue
				}
				if curline[0] == '#' {
					if strings.HasPrefix(curline, "# Category") {
						app.category = strings.Split(getcmntval(curline, "Category"), " ")
					}
					if strings.HasPrefix(curline, "# Icon") {
						app.icon = getcmntval(curline, "Icon")
					}
				}
			}
			applist = append(applist, app)
		}

	}
	return applist
}
