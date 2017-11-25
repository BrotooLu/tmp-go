package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("usage: sk-bundle infoPath srcPath destPath")
	}

	wd, _ := os.Getwd()
	log.Printf("exeloc: %s wd: %s\n", os.Args[0], wd)

	infoPath, _ := filepath.Abs(os.Args[1])
	infoFileInfo, err := os.Stat(infoPath)
	if err != nil {
		log.Fatalf("info file '%s' stat err: %s", infoPath, err)
	}

	infoFile, err := os.Open(infoPath)
	if err != nil {
		log.Fatalf("open info err: %s", err)
	}
	defer infoFile.Close()

	srcPath, _ := filepath.Abs(os.Args[2])
	srcFileInfo, err := os.Stat(srcPath)
	if err != nil {
		log.Fatalf("stat '%s' err: %s", srcPath, err)
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("err open src: %s", err)
	}
	defer srcFile.Close()

	destPath, err := filepath.Abs(os.Args[3])
	destFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("err open dest: %s", err)
	}
	defer destFile.Close()

	infoLen := infoFileInfo.Size()
	var lenWorld, emptyWorld [4]byte
	binary.BigEndian.PutUint32(lenWorld[:], uint32(infoLen))
	destFile.Write(lenWorld[:])
	binary.BigEndian.PutUint32(lenWorld[:], uint32(srcFileInfo.Size()))
	destFile.Write(lenWorld[:])
	destFile.Write(emptyWorld[:])
	destFile.Write(emptyWorld[:])

	buffer := make([]byte, infoLen)
	read, err := infoFile.Read(buffer)
	if err != nil {
		log.Fatalf("read info err: %s", err)
	}

	log.Printf("read info len: %d totoal: %d\n", read, infoLen)
	destFile.Write(buffer)

	buffer = make([]byte, 1)

	for {
		read, err = srcFile.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalf("read src err: %s", err)
		}

		if read == 0 {
			break
		}

		log.Printf("write %d %s to dest\n", read, buffer[:read])
		destFile.Write(buffer[:read])
	}
}
