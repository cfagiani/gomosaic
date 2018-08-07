package processor

import (
	"github.com/cfagiani/gomosaic"
	"sort"
)

type IndexProcessor interface {
	Process(oldIndex gomosaic.MosaicTiles, newIndex gomosaic.MosaicTiles) gomosaic.MosaicTiles
}

//performs a binary search of the sorted index for an entry with the filename specified
func find(name string, index gomosaic.MosaicTiles) *gomosaic.MosaicTile {
	i := sort.Search(len(index), func(i int) bool { return index[i].Filename == name })
	if i < len(index) && index[i].Filename == name {
		return &index[i]
	} else {
		return nil
	}
}
