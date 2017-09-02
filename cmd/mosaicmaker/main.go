package main

import (
	"os"
	"fmt"
)

//This command wil run the mosaic maker. It assumes that we have already computed an index to use for tiles.
func main() {
	if len(os.Args) < 4 {
		usage()
		os.Exit(1)
	}
	//TODO call mosaicmaker
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("go run mosaicindexer <imageDir>[,<imageDir>] size <outputDir>\n")
}
