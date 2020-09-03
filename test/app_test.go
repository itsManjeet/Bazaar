package main

import (
	"fmt"
	"os"
	"testing"
)

const recipie_dir = "/run/media/itsmanjeet/e9112081-7558-470f-8b6e-9dff8528fa2b/itsmanjeet/Work/releaxos/data/recipies/crux"

func TestFilePath(t *testing.T) {
	f, err := os.Open(recipie_dir)
	if err != nil {
		t.Fatal(err)
	}

	file_info, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		t.Fatal(err)
	}

	for _, file := range file_info {
		fmt.Println(file.Name())
	}

}
