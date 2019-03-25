package main

import (
	"io/ioutil"
	"log"
	"os"

	path2 "path"

	"strings"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	dryRun := true
	log.SetFlags(log.Lshortfile | log.Ltime)
	if len(os.Args) != 3 && len(os.Args) != 4 {
		log.Println("First argument - download directory, second argument torrents directory, third - dry run by default true")
		return
	}
	if len(os.Args) == 4 {
		dryRun = false
	}
	dir := os.Args[2]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	torrents := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".torrent") {
			tnm, err := tname(path2.Join(dir, f.Name()))
			if err == nil {
				torrents = append(torrents, tnm)
			} else {
				log.Println(err)
			}
		}
	}

	dir = os.Args[1]
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	totalf := len(files)
	totalo := 0
	for _, f := range files {
		found := false
		for i := range torrents {
			if torrents[i] == f.Name() {
				found = true
				break
			}
		}
		if !found {
			totalo = totalo + 1
			if dryRun {
				log.Printf("Orphan %s\n", f.Name())
			} else {
				log.Printf("Removed orphan %s\n", f.Name())
			}
		}
	}
	log.Printf("Orphans %d of %d\n", totalo, totalf)
}

func tname(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	mi, err := metainfo.Load(f)
	if err != nil {
		return "", err
	}

	var info metainfo.Info
	err = bencode.Unmarshal(mi.InfoBytes, &info)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}
