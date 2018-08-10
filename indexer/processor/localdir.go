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

const RecurseOption = "recurse"
const LocalKind = "local"

//Process will traverse a directory in a depth-first manner (if the option is set to recurse), looking for and analyzing any images. If the image is already in the
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
		if file.IsDir() && p.Source.Options == RecurseOption {
			processor := LocalProcessor{gomosaic.ImageSource{Options: RecurseOption, Path: filename, Kind: LocalKind}}
			newIndex = processor.Process(oldIndex, newIndex)
		} else if mosaicimages.IsSupportedImage(p.Source.Path, file) {
			existingTile := find(filename, oldIndex)
			if existingTile == nil {
				imageSegment, err := mosaicimages.AnalyzeImage(filename)
				if err == nil {
					//now add to index
					newIndex = append(newIndex,
						gomosaic.MosaicTile{Loc: "L", Filename: filename, AvgR: imageSegment.RVal, AvgG: imageSegment.GVal, AvgB: imageSegment.BVal})
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
