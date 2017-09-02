package mosaicmaker

import (
	"github.com/cfagiani/gomosaic/indexer"
	"log"
)

const (
	minIndexSize int = 100
)

func MakeMosaic(sourceImage string, indexDir string, gridSize uint32, tileSize uint32, outputFile string) {

	//read the index
	index := indexer.ReadIndex(indexDir)
	if len(index) < minIndexSize {
		log.Fatal("Index contains too few entries to generate a mosaic. Index  more tile images.")
	}


}
