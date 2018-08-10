package mosaicimages

import (
	"testing"
	"os"
)

func TestIsSupportedImage(t *testing.T) {
	cases := []struct {
		source    string
		supported bool
	}{
		{"../testdata/img1.png", true},
		{"../testdata/img3.jpg", true},
		{"../testdata/testconfig.json", false},
		{"../testdata/testindex.dat", false},
	}
	for _, c := range cases {
		file, err := os.Open(c.source)
		if err != nil {
			t.Errorf("Could not open file %v", err)
		}
		fileInfo, _ := file.Stat()
		isSupported := IsSupportedImage("../testdata", fileInfo)
		if isSupported != c.supported {
			t.Errorf("IsSupportedImage should have returned %v for %v but returned %v", c.supported, c.source, isSupported)
		}
		file.Close()
	}
}
