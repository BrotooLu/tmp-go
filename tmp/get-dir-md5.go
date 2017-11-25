package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

func computeMD5(path string, md5map map[string]string) error {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if info.IsDir() {
		files, err := file.Readdir(-1)
		if err != nil {
			return err
		}

		for _, f := range files {
			computeMD5(path+"/"+f.Name(), md5map)
		}
	} else {
		hash := md5.New()
		io.Copy(hash, file)
		md5map[path] = hex.EncodeToString(hash.Sum(nil))
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: get-dir-md5 dir")
	}

	path, _ := filepath.Abs(os.Args[1])
	log.Printf("path: %s", path)

	md5map := make(map[string]string)
	computeMD5(path, md5map)
	log.Printf("md5map: %v", md5map)
}
