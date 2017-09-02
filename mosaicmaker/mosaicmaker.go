package mosaicmaker

import (
	"github.com/cfagiani/gomosaic/indexer"
	"log"
)

const (
	//minimum number of entries in the index required to even attempt a mosaic
	minIndexSize int = 100
)

//Makes a new photomosaic of the sourceImage using the files referenced in the indexDir as a source. This method
//will divide up the source image into a grid and find the best match tile from the index to use in the output image.
func MakeMosaic(sourceImage string, indexDir string, gridSize uint32, tileSize uint32, outputFile string) {

	//read the index
	index := indexer.ReadIndex(indexDir)
	if len(index) < minIndexSize {
		log.Fatal("Index contains too few entries to generate a mosaic. Index  more tile images.")
	}

}
