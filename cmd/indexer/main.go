package main

import (
	"os"
	"fmt"
	"github.com/cfagiani/gomosaic/indexer"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	indexer.Index(os.Args[1], os.Args[2])
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("go run mosaicindexer <imageDir>[,<imageDir>] <outputDir>\n")
}
