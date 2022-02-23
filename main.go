package main

import (
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func handleFile(eventName string, NFS_PATH string) {
	fileInfo, err := os.Stat(eventName)
	if err != nil {
		panic("Couldn't determine the type of your file/folder")
	}
	var outputPath = NFS_PATH + "/" + path.Base(eventName)
	if !fileInfo.IsDir() {
		dirErr := os.MkdirAll(path.Dir(outputPath), os.ModePerm)
		if dirErr != nil {
			log.Panic(dirErr)
		}
	}

	renameError := os.Rename(eventName, outputPath)
	if renameError != nil {
		panic(err)
	}
}
func main() {
	NFS_PATH := os.Getenv("NFS_PATH")
	TARGET_PATH := os.Getenv("TARGET_PATH")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				fileIsDownloading := strings.Contains(path.Ext(event.Name), ".crdownload")
				if event.Op&fsnotify.Chmod == fsnotify.Chmod && !fileIsDownloading { // New File exists in the target folder By Host os means (copy/Create)

					log.Println("Copy Finished", event.Name)
					go handleFile(event.Name, NFS_PATH)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename && fileIsDownloading {
					time.Sleep(time.Millisecond * 200)
					log.Println("Download Finished", event.Name)
					var postRenameFile = strings.Replace(event.Name, ".crdownload", "", -1)
					go handleFile(postRenameFile, NFS_PATH)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(TARGET_PATH)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
