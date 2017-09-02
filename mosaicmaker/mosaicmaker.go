package mosaicmaker

import (
	"github.com/cfagiani/gomosaic/indexer"
	"log"
	"github.com/cfagiani/gomosaic/mosaicimages"

	"github.com/cfagiani/gomosaic"
)

const (
	//minimum number of entries in the index required to even attempt a mosaic
	minIndexSize int = 1
)

//Makes a new photomosaic of the sourceImage using the files referenced in the indexDir as a source. This method
//will divide up the source image into a grid and find the best match tile from the index to use in the output image.
func MakeMosaic(sourceImage string, indexDir string, gridSize int, tileSize int, outputFile string) {

	//read the index
	index := indexer.ReadIndex(indexDir)
	if len(index) < minIndexSize {
		log.Fatal("Index contains too few entries to generate a mosaic. Index  more tile images.")
	}
	segments, w, h, _ := mosaicimages.SegmentImage(sourceImage, gridSize)
	mosaic := make(map[gomosaic.ImageSegment]gomosaic.MosaicTile);
	for _, node := range segments {
		//TODO do this in parallel with goroutine/channels
		mosaic[node] = findBestTile(node, index)
	}

	//write final image
	outputImage := mosaicimages.CreateDrawableImage(tileSize, gridSize, w, h)
	for _, node := range segments {
		//TODO write image tile
		mosaicimages.WriteTileToImage(outputImage, mosaic[node], uint(tileSize), node.XMin, node.YMin)
	}
	//now write image to file
	mosaicimages.WriteImageToFile(outputImage, outputFile)

}

//Finds the tile with the closest match to the segment average color.
//TODO prevent duplicates
func findBestTile(segment gomosaic.ImageSegment, index gomosaic.MosaicTiles) gomosaic.MosaicTile {
	return index[0]
}
