package main

import (
	"os"
	"fmt"

)

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