package processor

import (
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/mosaicimages"
	"github.com/cfagiani/gomosaic/util"
	"io/ioutil"
	"log"
)

type LocalProcessor struct {
	Source gomosaic.ImageSource
}

//Processes a directory in a depth-first manner, looking for and analyzing any images. If the image is already in the
//index, the data will simply be copied to the new index without re-analyzing the image.
func (p LocalProcessor) Process(oldIndex gomosaic.MosaicTiles, newIndex gomosaic.MosaicTiles) gomosaic.MosaicTiles {
	files, err := ioutil.ReadDir(p.Source.Path)
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	log.Printf("Indexing %s\n", p.Source)
	for _, file := range files {
		filename := util.GetPath(p.Source.Path, file.Name())
		if file.IsDir() && p.Source.Options == "recurse" {
			processor := LocalProcessor{gomosaic.ImageSource{Options: "recurse", Path: filename, Kind: "local"}}
			newIndex = processor.Process(oldIndex, newIndex)
		} else if mosaicimages.IsSupportedImage(p.Source.Path, file) {
			existingTile := find(filename, oldIndex)
			if existingTile == nil {
				imageSegment, err := mosaicimages.AnalyzeImage(filename)
				if err == nil {
					//now add to index
					newIndex = append(newIndex,
						gomosaic.MosaicTile{Filename: filename, AvgR: imageSegment.RVal, AvgG: imageSegment.GVal, AvgB: imageSegment.BVal})
					count++
				}
			} else {
				newIndex = append(newIndex, *existingTile)
			}
		}
	}
	log.Printf("Added %d new files to index\n", count)
	return newIndex

}
