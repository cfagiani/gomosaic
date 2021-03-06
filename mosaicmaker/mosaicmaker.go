package mosaicmaker

import (
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/indexer"
	"github.com/cfagiani/gomosaic/mosaicimages"
	"github.com/cfagiani/gomosaic/util"
	"google.golang.org/api/photoslibrary/v1"
	"log"
	"math"
	"os"
	"errors"
)

const (
	//minimum number of entries in the index required to even attempt a mosaic
	minIndexSize = 50
	logInterval  = 10
)

//MakeMosaic makes a new photomosaic of the sourceImage using the files referenced in the indexDir as a source. This method
//will divide up the source image into a grid and find the best match tile from the index to use in the output image.
func MakeMosaic(sourceImage string, indexPath string, gridSize int, tileSize int, outputFile string, configFile string) error {
	var photoService *photoslibrary.Service
	var err error

	if gridSize <= 0 || tileSize <= 0 {
		return errors.New("gridSize and tileSize must be positive")
	}

	if len(configFile) > 0 {
		config, e := util.ReadConfig(configFile)
		if e != nil {
			log.Fatalf("Could not read configuration file: %v\n", e)
			os.Exit(1)
		}
		//TODO: get token from correct source
		photoService, err = util.GetPhotosService(config.GoogleClientId, config.GoogleClientSecret, "token.json")
		if err != nil {
			log.Fatalf("Could not initialize service client: %v\n", e)
			os.Exit(1)
		}
	}

	//read the index
	filename, exists := indexer.GetIndexFileName(indexPath)
	if !exists {
		log.Fatalf("Cannot produce a mosaic since index does not exist at %s", indexPath)
	}
	index := indexer.ReadIndex(filename)
	if len(index) < minIndexSize {
		log.Fatal("Index contains too few entries to generate a mosaic. Index  more tile images.")
	}
	log.Printf("Using index with %d entries", len(index))
	segments, w, h, _ := mosaicimages.SegmentImage(sourceImage, gridSize)
	mosaic := make(map[gomosaic.ImageSegment]gomosaic.MosaicTile)
	usedTiles := make(map[gomosaic.MosaicTile]bool)
	log.Println("Computing matches")
	for idx, node := range segments {
		//TODO do this in parallel with goroutine/channels
		//TODO come up with better findBestTile implementation
		//shouldn't be hard to improve on O(GI) where G is grid size and I is index size)
		mosaic[node] = findBestTile(node, index, usedTiles)
		if idx%logInterval == 0 {
			log.Printf("Tiles selected for %d segments", idx)
		}
	}

	log.Println("Assembling image")
	//write final image
	outputImage, _ := mosaicimages.CreateDrawableImage(tileSize, gridSize, w, h)
	for idx, node := range segments {
		x, y := projectToDestCoordinates(node, w, h, tileSize, gridSize)
		mosaicimages.WriteTileToImage(outputImage, mosaic[node], uint(tileSize), x, y, photoService)
		if idx%logInterval == 0 {
			log.Printf("Wrote %d tiles into destination image", idx)
		}
	}
	//now write image to file
	return mosaicimages.WriteImageToFile(outputImage, outputFile)

}

func projectToDestCoordinates(seg gomosaic.ImageSegment, w int, h int, tileSize int, gridSize int) (int, int) {
	tileX := seg.XMin / gridSize
	tileY := seg.YMin / gridSize
	return tileX * tileSize, tileY * tileSize
}

//Finds the tile with the closest match to the segment average color.
//TODO better duplicate handling
//TODO take a threshold as a param and exit once a match within that threshold is found (to avoid having to search the entire index each time)
func findBestTile(segment gomosaic.ImageSegment, index gomosaic.MosaicTiles, usedTiles map[gomosaic.MosaicTile]bool) gomosaic.MosaicTile {
	bestNode := index[0]
	bestDist := math.MaxFloat64
	for _, node := range index {
		if usedTiles[node] {
			continue
		}
		curDist := getDistance(segment, node)
		if curDist < bestDist {
			curDist = bestDist
			bestNode = node
		}
	}
	usedTiles[bestNode] = true
	return bestNode
}

//returns the "distance" between the tile and segment using the sum of squared distances calculation
func getDistance(segment gomosaic.ImageSegment, tile gomosaic.MosaicTile) float64 {
	return math.Pow(float64(segment.RVal-tile.AvgR), 2) + math.Pow(float64(segment.GVal-tile.AvgG), 2) + math.Pow(float64(segment.BVal-tile.AvgB), 2)
}
