package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type progressfunc func(int)

type downloader struct {
	current, total uint64
	status         string
	prgsfunc       progressfunc
}

func (dw *downloader) Write(p []byte) (int, error) {
	n := len(p)
	dw.current += uint64(n)
	dw.prgsfunc(int(float64(float64(dw.current)/float64(dw.total)) * 100))
	return n, nil
}

func (dw *downloader) download(filepath string, url string) error {
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

	fmt.Println(resp.Header.Get("Content-Length"))
	dw.total, err = strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return errors.New("error while getting file size")
	}
	defer resp.Body.Close()

	if _, err = io.Copy(out, io.TeeReader(resp.Body, dw)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")
	out.Close()

	if err = os.Rename(filepath+".part", filepath); err != nil {
		return err
	}

	return nil
}
