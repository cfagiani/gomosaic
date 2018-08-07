package indexer

import (
	"bufio"
	"fmt"
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/indexer/processor"
	"github.com/cfagiani/gomosaic/util"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"io/ioutil"
	"encoding/json"
)

const (
	//delimiter used in index file
	delimiter = ";"
	//name of index file
	idxname = "mosaicIndex.dat"
)

//Indexes all the readable images in the source dir (will recursively search all sources defined in the config file
//for images). For each image found, the average color values will be calculated and the results will be written to
//the dest file
func Index(configFile string, dest string) {

	file, e := ioutil.ReadFile(configFile)
	if e != nil {
		fmt.Printf("Could not read configuration file: %v\n", e)
		os.Exit(1)
	}
	var config gomosaic.Config
	json.Unmarshal(file, &config)

	//first read dat file if present
	log.Println("Reading file")
	oldIndex := ReadIndex(dest)
	log.Printf("Old index has %d entries\n", len(oldIndex))

	// TODO: use a goroutine for each directory?
	var newIndex gomosaic.MosaicTiles = make([]gomosaic.MosaicTile, 0, 100)
	for i := 0; i < len(config.Sources); i++ {
		sourceProcessor := getProcessor(config.Sources[i], config)
		newIndex = sourceProcessor.Process(oldIndex, newIndex)
	}
	oldIndex = nil // we don't need the old index anymore
	sort.Sort(newIndex)

	log.Printf("Writing new index with %d entries\n", len(newIndex))
	writeIndex(dest, newIndex)
}

//Reads an existing index and returns it as a MosaicTiles type. If the index does not exist, the MosaicTiles slice will
//be empty.
func ReadIndex(source string) gomosaic.MosaicTiles {
	var index = make([]gomosaic.MosaicTile, 0, 100)
	filename, exists := GetIndexFileName(source)
	if exists {
		f, err := os.Open(filename)
		util.CheckError(err, "Error opening file", true)
		//close file when block exits
		defer f.Close()
		r := bufio.NewReader(f)

		for {
			line, err := r.ReadString(10) // 0x0A separator = newline
			if err == nil {
				index = append(index, createNodeFromLine(line))
			}
			if err == io.EOF {
				break
			}
		}
	}
	return index
}

//returns the filename that should be used for the index along with a flag indicating if the file exists
func GetIndexFileName(source string) (string, bool) {
	filename := source
	if fileInfo, err := os.Stat(source); !os.IsNotExist(err) {
		exists := true
		if fileInfo.IsDir() {
			//if the source existed and was a directory, add the default index name
			filename = source + string(os.PathSeparator) + idxname
			//now see if that exists
			_, err := os.Stat(filename)
			exists = !os.IsNotExist(err)
		}
		return filename, exists
	} else {
		return filename, false
	}
}

//Prints the entire index.
func printIndex(index []gomosaic.MosaicTile) {
	for _, node := range index {
		fmt.Println(node.ToString())
	}
}

//Writes the index file to the destDir.
func writeIndex(dest string, index gomosaic.MosaicTiles) {
	filename, _ := GetIndexFileName(dest)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	util.CheckError(err, "error opening file", true)
	// close file when block exits
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, node := range index {
		if node.Filename != "" {
			fmt.Fprintf(w, "%s\n", node.ToString())
		}
	}
}

func getProcessor(source gomosaic.ImageSource, config gomosaic.Config) processor.IndexProcessor {
	if source.Kind == "local" {
		return processor.LocalProcessor{Source: source}
	} else {
		return processor.GooglePhotosProcess{Source: source, Config: config}
	}
}

//Parses a line from the index and uses it to initialize a new MosaicTile
func createNodeFromLine(line string) gomosaic.MosaicTile {
	// construct node
	parts := strings.Split(line, delimiter)
	return gomosaic.MosaicTile{Filename: parts[0], AvgR: util.GetInt32(parts[1]), AvgG: util.GetInt32(parts[2]), AvgB: util.GetInt32(parts[3])}

}
