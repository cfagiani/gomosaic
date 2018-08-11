package mosaicimages

import (
	"os"
	"testing"
)

//TestIsSupportedImage will ensure that the IsSupportedImage method returns true for supported files and
//false for non-supported files.
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

//TestResizeImage validates that the ResizeImage returns an Image struct with the desired bounds as long as the source is
//a supported image file.
func TestResizeImage(t *testing.T) {
	cases := []struct {
		source      string
		width       uint
		height      uint
		expectError bool
	}{
		{"../testdata/img1.png", 10, 10, false},
		{"../testdata/img3.jpg", 200, 300, false},
		{"../testdata/testindex.dat", 10, 10, true},
		{"../testdata/notThere", 10, 10, true},
	}
	for _, c := range cases {
		img, err := ResizeImage(c.source, c.height, c.width)
		if err != nil && !c.expectError {
			t.Errorf("ResizeImage returned an unexpected error for %v", c.source)
		} else if err == nil {
			if c.expectError {
				t.Errorf("ResizeImage should have returned an error for %v but did not", c.source)
			} else {
				//now verify the resize is as desired
				if img.Bounds().Size().X != int(c.width) || img.Bounds().Size().Y != int(c.height) {
					t.Errorf("Image bound not correct after ResizeImage for %v. Wanted %vx%v but got %vx%v", c.source,
						c.width, c.height, img.Bounds().Size().X, img.Bounds().Size().Y)
				}
			}
		}
	}
}

//TestCreateDrawableImage validates that the CreateDrawableImage returns an image with the correct bounds.
func TestCreateDrawableImage(t *testing.T) {
	cases := []struct {
		t           int
		g           int
		w           int
		h           int
		expectError bool
	}{
		{5, 10, 100, 100, false},
		{5, 10, -1, 100, true},
		{5, 10, 0, 0, true},
		{5, 100, 10, 10, true},
		{0, 10, 100, 100, true},
	}
	for _, c := range cases {
		img, err := CreateDrawableImage(c.t, c.g, c.w, c.h)
		if err != nil && !c.expectError {
			t.Errorf("Expected an error from CreateDrawableImage for %v but did not get one", c)
		} else if err == nil {
			if c.expectError {
				t.Errorf("Did not get an error from CreateDrawableImage but should have for %v", c)
			}else{
				if img.Bounds().Size().X != (c.w/c.g)*c.t ||  img.Bounds().Size().Y != (c.h/c.g)*c.t{
					t.Errorf("CreateDrawableImage returned image of wrong dimensions for %v", c)
				}
			}
		}
	}
}
