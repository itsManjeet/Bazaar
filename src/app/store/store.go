package store

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gotk3/gotk3/gdk"
)

// App interface
type App interface {
	ID() string
	Name() string
	Version() string
	Release() string
	Description() string
	Depends() []string
	PackageFile() string
	Url() string
	Categories() []string
	License() string
	Store() string
	Icon(size uint) *gdk.Pixbuf
}

// Store interface
type Store interface {
	ID() string
	ListApps() ([]App, error)
	GetApp(appname string) (App, error)
	Depends(appname string) ([]App, error)
	Install(appname string, prg ProgFunc) error
	IsInstalled(appname string) bool
}

// ProgFunc progress function type
type ProgFunc func(int)

// Downloader downloader struct
type Downloader struct {
	Current, Total   float64
	Status           string
	ProgressFunction ProgFunc
}

func (dw *Downloader) Write(p []byte) (int, error) {
	n := len(p)
	dw.Current += float64(n)
	dw.ProgressFunction(int(dw.Current/dw.Total) * 100)
	return n, nil
}

// Download download file
func (dw *Downloader) Download(filepath, url string) error {
	out, err := os.Create(filepath + ".part")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	log.Println(resp.Header.Get("Content-Length"))
	dw.Total, err = strconv.ParseFloat(resp.Header.Get("Content-Length"), 64)
	if err != nil {
		return errors.New(err.Error())
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, io.TeeReader(resp.Body, dw)); err != nil {
		out.Close()
		return err
	}

	out.Close()
	return os.Rename(filepath+".part", filepath)
}

// Extractor extract tarfile
func Extractor(tarfile, destdir string) error {
	tr, err := os.Open(tarfile)
	if err != nil {
		log.Println("Extractor()->os.Open()", tarfile, err)
		return err
	}

	uncstr, err := gzip.NewReader(tr)
	if err != nil {
		log.Println("Exractor()->gzip.NewReader()", tarfile, err)
		return err
	}

	tarrdr := tar.NewReader(uncstr)

	for true {
		header, err := tarrdr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Extractor()->tarrdr.Next() failed:", err)
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(destdir, header.Name), os.FileMode(header.Mode)); err != nil {
				log.Println("Extractor()->MkdirAll()  [TypeDir] failed: ", err)
				return err
			}

		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				log.Println("Extractor()->Create() [TypeReg] failed: ", err)
				return err
			}

			if _, err := io.Copy(outFile, tarrdr); err != nil {
				log.Println("Extractor()->Copy() [TypeReg] failed:", err)
			}
			outFile.Close()

		default:
			log.Println("Extractor(): [UnknownType]: ", header.Typeflag, "in", header.Name)
			return err
		}
	}

	return nil
}
