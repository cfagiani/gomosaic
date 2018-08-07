package mosaicimages

import (
	"errors"
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/util"
	"github.com/nfnt/resize"
	"github.com/utahta/go-openuri"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"strings"
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
	f, err := os.Open(util.GetPath(dirName, file.Name()))
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
func resizeImage(inputFile string, height uint, width uint) image.Image {
	file, err := os.Open(inputFile)
	defer file.Close()
	util.CheckError(err, "Could not read input file", true)

	// decode jpeg into image.Image
	img, _, err := image.Decode(file)
	util.CheckError(err, "Could not decode image", true)
	return resize.Resize(width, height, img, resize.Lanczos3)
}

//resizes an input image and writes a new file with the result
func GenerateTileFile(inputFile string, outputFile string, height uint, width uint) {
	m := resizeImage(inputFile, height, width)
	WriteImageToFile(m, outputFile)
}

//Creates a new Image using the dimensions passed in
func CreateDrawableImage(tileSize int, gridSize int, sourceWidth int, sourceHeight int) draw.Image {
	return image.NewRGBA(image.Rect(0, 0, (sourceWidth/gridSize)*tileSize, (sourceHeight/gridSize)*tileSize))
}

func WriteTileToImage(img draw.Image, tile gomosaic.MosaicTile, tileSize uint, startX int, startY int) {
	tileImage := resizeImage(tile.Filename, tileSize, tileSize)
	destRec := image.Rect(startX, startY, startX+int(tileSize), startY+int(tileSize))
	draw.FloydSteinberg.Draw(img, destRec.Bounds(), tileImage,
		image.Point{tileImage.Bounds().Min.X, tileImage.Bounds().Min.Y})
}

func WriteImageToFile(img image.Image, outputFile string) {
	out, err := os.Create(outputFile)
	util.CheckError(err, "Could not open output file", true)
	defer out.Close()
	// write new image to file
	jpeg.Encode(out, img, nil)
}

//Divides a source image up into square segments of the specified size and returns an array of ImageSegments. If the
//image cannot be processed, an error is returned.
func SegmentImage(sourceImage string, segmentSize int) ([]gomosaic.ImageSegment, int, int, error) {
	file, err := os.Open(sourceImage)
	if !util.CheckError(err, "Could not process image", false) {
		defer file.Close()
		img, _, err := image.Decode(file)
		util.CheckError(err, "Could not process image", true)
		bounds := img.Bounds()
		//TODO: need to handle non-square images better
		var segments []gomosaic.ImageSegment = make([]gomosaic.ImageSegment, 0, 100)
		for y := bounds.Min.Y; y < bounds.Max.Y; y += segmentSize {
			for x := bounds.Min.X; x < bounds.Max.X; x += segmentSize {
				segments = append(segments, analyzeImageSegment(img, x, y, x+segmentSize, y+segmentSize))
			}
		}
		return segments, bounds.Max.X - bounds.Min.X, bounds.Max.Y - bounds.Min.Y, nil
	} else {
		return make([]gomosaic.ImageSegment, 0, 0), 0, 0, errors.New("Could not analyze image")
	}
}

//Analyzes an entire image and returns an ImageSegment with the result. If the image cannot be decoded, an error is
//returned.
func AnalyzeImage(filename string) (gomosaic.ImageSegment, error) {
	file, err := openuri.Open(filename)

	if !util.CheckError(err, "Could not process image", false) {
		defer file.Close()
		img, _, err := image.Decode(file)
		util.CheckError(err, "Could not process image", true)
		bounds := img.Bounds()
		return analyzeImageSegment(img, bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y), nil
	} else {
		return gomosaic.ImageSegment{0, 0, 0, 0, 0, 0, 0}, errors.New("Could not analyze image")
	}
}

//Analyzes a segment of an image, returning an ImageSegment struct with the result.
func analyzeImageSegment(img image.Image, xMin int, yMin int, xMax int, yMax int) gomosaic.ImageSegment {
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
	return gomosaic.ImageSegment{xMin, yMin, xMax, yMax,
		rTotal / pixelCount, gTotal / pixelCount, bTotal / pixelCount}
}
