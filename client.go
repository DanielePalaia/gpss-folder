package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: you need to specify 1 parameter the directory path")
	}
	/* Read the directory field path */
	path := os.Args[1]

	/* Remove last / if present*/
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}

	/* Run the directory management object */
	folderEngine := makeFolderEngine(path)
	folderEngine.createOperationalFolders()
	folderEngine.folderListen()

}
