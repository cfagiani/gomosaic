package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sort"
	"errors"
)

const (
	_delimiter_ string = ";"
	_idxname_   string = "mosaicIndex.dat"
)

var magicNumbers = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

func main() {
	if len(os.Args) < 4 {
		usage()
		os.Exit(1)
	}

	//first read dat file if present
	log.Println("Reading file")
	oldIndex := readIndex(os.Args[3])
	log.Printf("Old index has %d entries\n", len(oldIndex))

	//now recurse through the directories
	dirs := strings.Split(os.Args[1], ",")
	sort.Strings(dirs)

	// TODO: use a goroutine for each directory?
	var newIndex MosaicTiles = make([]MosaicTile, 0, 100)
	for i := 0; i < len(dirs); i++ {
		newIndex = processDirectory(dirs[i], os.Args[3], getInt(os.Args[2]), oldIndex, newIndex)
	}
	oldIndex = nil // we don't need the old index anymore
	sort.Sort(newIndex)

	log.Printf("Writing new index with %d entries\n", len(newIndex))
	writeIndex(os.Args[3], newIndex)
}

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("go run mosaicindexer <imageDir>[,<imageDir>] size <outputDir>\n")
}

func processDirectory(dirName string, outputDir string, desiredSize uint32, oldIndex MosaicTiles, newIndex MosaicTiles) MosaicTiles {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", dirName, file.Name())
		if file.IsDir() {
			newIndex = processDirectory(filename, outputDir, desiredSize, oldIndex, newIndex)
		} else if isSupportedImage(dirName, file) {
			existingTile := find(filename, oldIndex)
			if existingTile == nil {
				rval, gval, bval, err := analyze(filename)
				if err == nil {
					//now add to index
					newIndex = append(newIndex,
						MosaicTile{filename, rval, gval, bval})
				}
			} else {
				newIndex = append(newIndex, *existingTile)
			}
		}
	}
	return newIndex
}

func analyze(filename string) (uint32, uint32, uint32, error) {
	file, err := os.Open(filename)

	if !checkError(err, "Could not process image", false) {

		defer file.Close()
		m, _, err := image.Decode(file)
		checkError(err, "Could not process image", true)
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

func printIndex(index []MosaicTile) {
	for _, node := range index {
		fmt.Println(node.ToString())
	}
}

func getAbsolutePath(dir string, file string) string {
	return dir + string(os.PathSeparator) + file
}

func checkError(err error, msg string, isFatal bool) bool {
	if err != nil {
		if isFatal {
			log.Fatal(msg, err)
		} else {
			log.Println(msg)
		}
		return true
	}
	return false
}

func writeIndex(destDir string, index MosaicTiles) {
	f, err := os.OpenFile(getAbsolutePath(destDir, _idxname_), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	checkError(err, "error opening file", true)
	// close file when block exits
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, node := range index {
		if node.filename != "" {
			fmt.Fprintf(w, "%s\n", node.ToString())
		}
	}
}

//checks if a file is a supported image by looking at the first few bytes to see if its in our magicNumber table
//while we could use the Decode method from images, we don't need/want to read the whole file right now
func isSupportedImage(dirName string, file os.FileInfo) bool {
	f, err := os.Open(getAbsolutePath(dirName, file.Name()))
	defer f.Close()
	if !checkError(err, "error opening file", false) {
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

//performs a binary search of the sorted index for an entry with the filename specified
func find(name string, index MosaicTiles) *MosaicTile {
	i := sort.Search(len(index), func(i int) bool { return index[i].filename == name })
	if i < len(index) && index[i].filename == name {
		return &index[i]
	} else {
		return nil
	}
}

func readIndex(src string) MosaicTiles {
	var index = make([]MosaicTile, 0, 100)
	filename := src + string(os.PathSeparator) + _idxname_
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		f, err := os.Open(filename)
		checkError(err, "Error opening file", true)
		//close file when block exits
		defer f.Close()
		r := bufio.NewReader(f)

		for {
			line, err := r.ReadString(10) // 0x0A separator = newline
			if err == nil {
				//insertNodeFromLine(line, indexRoot)
				index = append(index, createNodeFromLine(line))
			}
			if err == io.EOF {
				break
			}
		}
	}
	return index
}

func createNodeFromLine(line string) MosaicTile {
	// construct node
	parts := strings.Split(line, _delimiter_)
	return MosaicTile{parts[0], getInt(parts[1]), getInt(parts[2]), getInt(parts[3])}
}

// Converts a string to an integer, eating any errors
func getInt(s string) uint32 {
	// TODO: actually handle the error
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint32(i)
	}
	return 0
}

type MosaicTile struct {
	filename string
	avgR     uint32
	avgG     uint32
	avgB     uint32
}

func (t MosaicTile) ToString() string {
	return fmt.Sprintf("%s;%d;%d;%d", t.filename, t.avgR, t.avgG, t.avgB)
}

//define a type so we can implement Sort interface
type MosaicTiles []MosaicTile

func (slice MosaicTiles) Len() int {
	return len(slice)
}

func (slice MosaicTiles) Less(i, j int) bool {
	return slice[i].filename < slice[j].filename
}

func (slice MosaicTiles) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
