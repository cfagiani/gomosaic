package indexer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"github.com/cfagiani/gomosaic/util"
	"github.com/cfagiani/gomosaic/mosaicimages"
)

const (
	//delimiter used in index file
	delimiter string = ";"
	//name of index file
	idxname string = "mosaicIndex.dat"
)

//Indexes all the readable images in the source dir (will recursively search all subdirs of the directories passed in
//for images). For each image found, the average color values will be calculated and the results will be written to a
//file in the destDir directory.
func Index(sourceDirs string, destDir string) {
	//first read dat file if present
	log.Println("Reading file")
	oldIndex := ReadIndex(destDir)
	log.Printf("Old index has %d entries\n", len(oldIndex))

	//now recurse through the directories
	dirs := strings.Split(sourceDirs, ",")
	sort.Strings(dirs)

	// TODO: use a goroutine for each directory?
	var newIndex MosaicTiles = make([]MosaicTile, 0, 100)
	for i := 0; i < len(dirs); i++ {
		newIndex = processDirectory(dirs[i], oldIndex, newIndex)
	}
	oldIndex = nil // we don't need the old index anymore
	sort.Sort(newIndex)

	log.Printf("Writing new index with %d entries\n", len(newIndex))
	writeIndex(destDir, newIndex)
}

//Processes a directory in a depth-first manner, looking for and analyzing any images. If the image is already in the
//index, the data will simply be copied to the new index without re-analyzing the image.
func processDirectory(dirName string, oldIndex MosaicTiles, newIndex MosaicTiles) MosaicTiles {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", dirName, file.Name())
		if file.IsDir() {
			newIndex = processDirectory(filename, oldIndex, newIndex)
		} else if mosaicimages.IsSupportedImage(dirName, file) {
			existingTile := find(filename, oldIndex)
			if existingTile == nil {
				imageSegment, err := mosaicimages.AnalyzeImage(filename)
				if err == nil {
					//now add to index
					newIndex = append(newIndex,
						MosaicTile{filename, imageSegment.RVal, imageSegment.GVal, imageSegment.BVal})
				}
			} else {
				newIndex = append(newIndex, *existingTile)
			}
		}
	}
	return newIndex
}

//Prints the entire index.
func printIndex(index []MosaicTile) {
	for _, node := range index {
		fmt.Println(node.ToString())
	}
}

//Writes the index file to the destDir.
func writeIndex(destDir string, index MosaicTiles) {
	f, err := os.OpenFile(util.GetAbsolutePath(destDir, idxname), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	util.CheckError(err, "error opening file", true)
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

//performs a binary search of the sorted index for an entry with the filename specified
func find(name string, index MosaicTiles) *MosaicTile {
	i := sort.Search(len(index), func(i int) bool { return index[i].filename == name })
	if i < len(index) && index[i].filename == name {
		return &index[i]
	} else {
		return nil
	}
}

//Reads an existing index and returns it as a MosaicTiles type. If the index does not exist, the MosaicTiles slice will
//be empty.
func ReadIndex(sourceDir string) MosaicTiles {
	var index = make([]MosaicTile, 0, 100)
	filename := sourceDir + string(os.PathSeparator) + idxname
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		f, err := os.Open(filename)
		util.CheckError(err, "Error opening file", true)
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

//Parses a line from the index and uses it to initialize a new MosaicTile
func createNodeFromLine(line string) MosaicTile {
	// construct node
	parts := strings.Split(line, delimiter)
	return MosaicTile{parts[0], util.GetInt(parts[1]), util.GetInt(parts[2]), util.GetInt(parts[3])}
}

//Type representing a tile that can be used in a mosaic
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
