package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/radovskyb/watcher"
)

type folderEngine struct {
	path string
}

func makeFolderEngine(path string) *folderEngine {

	engine := new(folderEngine)
	engine.path = path
	return engine
}

func (engine *folderEngine) createOperationalFolders() {

	createDirectory(engine.path + "/staging")
	createDirectory(engine.path + "/processing")
	createDirectory(engine.path + "/completed")
}

func (engine *folderEngine) folderListen() {

	w := watcher.New()

	// Only notify create events.
	w.FilterOps(watcher.Create)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event.Path) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(engine.path); err != nil {
		log.Fatalln(err)
	}

	// Trigger 2 events after watcher started.
	go func() {
		w.Wait()
		w.TriggerEvent(watcher.Create, nil)
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

}

func createDirectory(dirName string) bool {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			panic(err)
		}
		return true
	}

	if src.Mode().IsRegular() {
		fmt.Println(dirName, "already exist as a file!")
		return false
	}

	return false
}
