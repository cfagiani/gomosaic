package mosaicimages

import (
	"os"
	"strings"
	"github.com/nfnt/resize"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"image"
	"errors"
	"github.com/cfagiani/gomosaic/util"
)

var magicNumbers = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

//checks if a file is a supported image by looking at the first few bytes to see if its in our magicNumber table
//while we could use the Decode method from images, we don't need/want to read the whole file right now
func IsSupportedImage(dirName string, file os.FileInfo) bool {
	f, err := os.Open(util.GetAbsolutePath(dirName, file.Name()))
	defer f.Close()
	if !util.CheckError(err, "error opening file", false) {
		var header = make([]byte, 36)
		f.Read(header) //we don't care about the error here since we'll just skip it
		headerStr := string(header)
		for magic := range magicNumbers {
			if strings.HasPrefix(headerStr, magic) {
				return true
			}
		}
	}
	return false
}

// use height or width of 0 to preserve aspect ratio
func Resize(inputFile string, outputFile string, height uint, width uint) {
	// open "test.jpg"
	file, err := os.Open(inputFile)
	defer file.Close()
	util.CheckError(err, "Could not read input file", true)

	// decode jpeg into image.Image
	img, _, err := image.Decode(file)
	util.CheckError(err, "Could not decode image", true)
	file.Close()

	m := resize.Resize(width, height, img, resize.Lanczos3)

	out, err := os.Create(outputFile)
	util.CheckError(err, "Could not open output file", true)
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

//Analyzes an entire image and returns an ImageSegment with the result. If the image cannot be decoded, an error is
//returned.
func AnalyzeImage(filename string) (ImageSegment, error) {
	file, err := os.Open(filename)

	if !util.CheckError(err, "Could not process image", false) {

		defer file.Close()
		img, _, err := image.Decode(file)
		util.CheckError(err, "Could not process image", true)
		bounds := img.Bounds()
		return analyzeImageSegment(img, bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y), nil
	} else {
		return ImageSegment{0, 0, 0, 0, 0, 0, 0}, errors.New("Could not analyze image")
	}
}

//Analyzes a segment of an image, returning an ImageSegment struct with the result.
func analyzeImageSegment(img image.Image, xMin int, yMin int, xMax int, yMax int) ImageSegment {
	var rTotal, gTotal, bTotal, pixelCount uint32 = 0, 0, 0, 0

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			rTotal += r
			gTotal += g
			bTotal += b
			pixelCount++
		}
	}
	return ImageSegment{xMin, yMin, xMax, yMax,
		rTotal / pixelCount, gTotal / pixelCount, bTotal / pixelCount}
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
