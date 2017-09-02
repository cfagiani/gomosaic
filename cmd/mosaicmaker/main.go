package main

import (
	"os"
	"fmt"
	"strconv"
	"github.com/cfagiani/gomosaic/mosaicmaker"
)

//This command wil run the mosaic maker. It assumes that we have already computed an index to use for tiles.
func main() {
	if len(os.Args) < 6 {
		usage()
		os.Exit(1)
	}
	gridSize, _ := strconv.Atoi(os.Args[3])
	tileSize, _ := strconv.Atoi(os.Args[4])
	mosaicmaker.MakeMosaic(os.Args[1], os.Args[2], gridSize, tileSize, os.Args[5])
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("go run cmd/mosaicmaker/main.go <sourceImage> <indexDir> <gridSize> <tileSize> <outputFile>\n")
}
