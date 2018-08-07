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
	indexer.Index(os.Args[1], os.Args[2])
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("go run cmd/indexer/main.go <configFile> <indexFile>\n")
}
