package main

import (
	"bufio"
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
				fmt.Println("New file detected Processing: " + event.Path)
				go engine.ProcessFile(event.Path)

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(engine.path + "/staging"); err != nil {
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

func (engine *folderEngine) ProcessFile(path string) {

	_, fileName := returnFileFromPath(path)
	fmt.Println(fileName)

	present := returnAndCompareAllFiles(fileName, engine.path)
	if present == true {
		fmt.Println("File: " + path + " already processed and present in completed, files are also equals Skip")
	} else {
		fmt.Println("File not yet processed, Processing: " + path)
		// Copy the file in the processing folder
		if _, err := copy(path, engine.path+"/processing/"+fileName); err != nil {
			fmt.Println("error is: ", err.Error())
			return

		}
		// Process the file
		engine.SendFileKafka(engine.path + "/processing/" + fileName)

		// Move the file from processing to completed
		err := os.Rename(engine.path+"/processing/"+fileName, engine.path+"/completed/"+fileName)
		if err != nil {
			log.Fatal(err)
		}

	}

}

// Read line by line
func (engine *folderEngine) SendFileKafka(path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
