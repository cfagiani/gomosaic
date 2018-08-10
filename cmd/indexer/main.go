package main

import (
	"fmt"
	"github.com/cfagiani/gomosaic/indexer"
	"os"
)

//This command will run the mosaic indexer on all the directories passed in via the command line. The index will be
//written to the output directory as specified on the command line.
func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	err := indexer.Index(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("Error while indexing %v", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("indexer <configFile> <indexFile>\n")
}
