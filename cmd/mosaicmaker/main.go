package main

import (
	"fmt"
	"github.com/cfagiani/gomosaic/mosaicmaker"
	"os"
	"strconv"
)

//This command wil run the mosaic maker. It assumes that we have already computed an index to use for tiles.
func main() {
	if len(os.Args) < 6 {
		usage()
		os.Exit(1)
	}
	gridSize, _ := strconv.Atoi(os.Args[3])
	tileSize, _ := strconv.Atoi(os.Args[4])
	configFile := ""
	if len(os.Args) == 7 {
		configFile = os.Args[6]
	}
	mosaicmaker.MakeMosaic(os.Args[1], os.Args[2], gridSize, tileSize, os.Args[5], configFile)
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("go run cmd/mosaicmaker/main.go <sourceImage> <index> <gridSize> <tileSize> <outputFile>\n")
}
