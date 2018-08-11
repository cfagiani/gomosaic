package mosaicmaker

import (
	"testing"
	"github.com/cfagiani/gomosaic"
)

func TestProjectToDestCoordinates(t *testing.T) {
	cases := []struct {
		seg       gomosaic.ImageSegment
		w         int
		h         int
		tileSize  int
		gridSize  int
		expectedX int
		expectedY int
	}{
		{gomosaic.ImageSegment{0, 0, 100, 100, 1, 1, 1}, 200, 200, 50, 50, 0, 0},
		{gomosaic.ImageSegment{100, 0, 100, 100, 1, 1, 1}, 200, 200, 50, 50, 100, 0},
		{gomosaic.ImageSegment{100, 0, 100, 100, 1, 1, 1}, 200, 200, 10, 50, 20, 0},
		{gomosaic.ImageSegment{100, 0, 100, 100, 1, 1, 1}, 200, 200, 10, 5, 200, 0},
		{gomosaic.ImageSegment{100, 100, 200, 200, 1, 1, 1}, 200, 200, 10, 5, 200, 200},
	}
	for _, c := range cases {
		x, y := projectToDestCoordinates(c.seg, c.w, c.h, c.tileSize, c.gridSize)
		if x != c.expectedX || y != c.expectedY {
			t.Errorf("projectToDestCoordinates returned %v,%v when %v,%v was expected", x, y, c.expectedX, c.expectedY)
		}
	}
}
