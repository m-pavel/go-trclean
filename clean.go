package main

import (
	"io/ioutil"
	"log"
	"os"

	path2 "path"

	"strings"

	"fmt"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	dryRun := true
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("First argument - download directory, second argument torrents directory, third - dry run by default true")
		return
	}
	if len(os.Args) == 4 {
		dryRun = false
	}
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	torrents := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".torrent") {
			tnm, err := tname(path2.Join(os.Args[2], f.Name()))
			if err != nil {
				torrents = append(torrents, tnm)
			} else {
				fmt.Println(err)
			}
		}
	}

	files, err = ioutil.ReadDir(os.Args[1])
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
				fmt.Printf("Orphan %s\n", f.Name())
			} else {
				fmt.Printf("Removed orphan %s\n", f.Name())
			}
		}
	}
	fmt.Printf("Orphans %d of %d\n")
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
