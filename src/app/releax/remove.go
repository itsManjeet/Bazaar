package releax

import (
	"bufio"
	"log"
	"os"
	"path"
)

// Remove appfile
func (s Store) Remove(appname string) error {
	file, err := os.Open(path.Join(s.DataDir, appname, "files"))
	if err != nil {
		log.Println("Remove() failed:", err)
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		curfile := path.Join("/", sc.Text())
		if exists(curfile) {
			log.Println("cleaning: ", curfile)
			err := os.Remove(curfile)
			if err != nil {
				log.Println("Remove() warning unable to remove", curfile, ": ", err)
			}
		}
	}

	return nil
}
