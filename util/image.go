package util

import (
	"os"
	"strings"
	"github.com/nfnt/resize"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"image"
	"errors"
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
	f, err := os.Open(GetAbsolutePath(dirName, file.Name()))
	defer f.Close()
	if !CheckError(err, "error opening file", false) {
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
	CheckError(err, "Could not read input file", true)

	// decode jpeg into image.Image
	img, _, err := image.Decode(file)
	CheckError(err, "Could not decode image", true)
	file.Close()

	m := resize.Resize(width, height, img, resize.Lanczos3)

	out, err := os.Create(outputFile)
	CheckError(err, "Could not open output file", true)
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

func AnalyzeImage(filename string) (uint32, uint32, uint32, error) {
	file, err := os.Open(filename)

	if !CheckError(err, "Could not process image", false) {

		defer file.Close()
		m, _, err := image.Decode(file)
		CheckError(err, "Could not process image", true)
		bounds := m.Bounds()
		var rTotal, gTotal, bTotal, pixelCount uint32 = 0, 0, 0, 0

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				// A color's RGBA method returns values in the range [0, 65535].
				rTotal += r
				gTotal += g
				bTotal += b
				pixelCount++
			}
		}
		return rTotal / pixelCount, gTotal / pixelCount, bTotal / pixelCount, nil
	} else {
		return 0, 0, 0, errors.New("Could not analyze image")
	}

}