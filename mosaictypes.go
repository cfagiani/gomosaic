package gomosaic

import (
	"fmt"
)

type Config struct {
	GoogleClientId     string
	GoogleClientSecret string
	Sources            []ImageSource
}

type ImageSource struct {
	Kind    string
	Path    string
	Options string
}

//Type representing a tile that can be used in a mosaic
type MosaicTile struct {
	Loc      string
	Filename string
	AvgR     uint32
	AvgG     uint32
	AvgB     uint32
}

func (t MosaicTile) ToString() string {
	return fmt.Sprintf("%s;%s;%d;%d;%d", t.Loc, t.Filename, t.AvgR, t.AvgG, t.AvgB)
}

//define a type so we can implement Sort interface
type MosaicTiles []MosaicTile

func (slice MosaicTiles) Len() int {
	return len(slice)
}

func (slice MosaicTiles) Less(i, j int) bool {
	return slice[i].Filename < slice[j].Filename
}

func (slice MosaicTiles) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//type representing the average color values of a segment of an image defined by the min/max X/Y coordinates
type ImageSegment struct {
	XMin int
	YMin int
	XMax int
	YMax int
	RVal uint32
	GVal uint32
	BVal uint32
}

func (t ImageSegment) ToString() string {
	return fmt.Sprintf("(%d,%d) to (%d,%d): %d;%d;%d", t.XMin, t.YMin, t.XMax, t.YMax, t.RVal, t.GVal, t.BVal)
}
