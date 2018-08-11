package indexer

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/indexer/processor"
	"github.com/cfagiani/gomosaic/util"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	//delimiter used in index file
	delimiter = ";"
	//name of index file
	idxname = "mosaicIndex.dat"
)

//Index will process all the readable images in the sources defined in the configuration file. For each image found,
//the average color values will be calculated and the results will be written to the dest file so it can be used in
//subsequent mosaic creations.
func Index(configFile string, dest string) error {

	config, e := util.ReadConfig(configFile)
	if e != nil {
		log.Fatalf("Could not read configuration file: %v\n", e)
		return e
	}

	//first read existing index file if present
	log.Println("Reading file")
	oldIndex := ReadIndex(dest)
	log.Printf("Old index has %d entries\n", len(oldIndex))

	// TODO: use a goroutine for each source?
	var newIndex gomosaic.MosaicTiles = make([]gomosaic.MosaicTile, 0, 100)
	for i := 0; i < len(config.Sources); i++ {
		sourceProcessor := getProcessor(config.Sources[i], config)
		if sourceProcessor == nil {
			log.Println("Skipping source.")
		}
		newIndex = sourceProcessor.Process(oldIndex, newIndex)
	}
	oldIndex = nil // we don't need the old index anymore
	sort.Sort(newIndex)

	log.Printf("Writing new index with %d entries\n", len(newIndex))
	writeIndex(dest, newIndex)
	return nil
}

//ReadIndex reads an existing index and returns it as a MosaicTiles type. If the index does not exist, the MosaicTiles slice will
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
				tile, lineErr := createNodeFromLine(line)
				if lineErr == nil {
					index = append(index, *tile)
				} else {
					log.Println("Ignoring invalid index line")
				}
			}
			if err == io.EOF {
				break
			}
		}
	}
	return index
}

//GetIndexFileName returns the filename that should be used for the index along with a flag indicating if the file exists
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

//writeIndex writes the index file to the destDir.
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

//getProcessor will return an instance of a type that implements the IndexProcessor interface.
func getProcessor(source gomosaic.ImageSource, config gomosaic.Config) processor.IndexProcessor {
	if source.Kind == processor.LocalKind {
		return processor.LocalProcessor{Source: source}
	} else if source.Kind == processor.GoogleKind {
		return processor.GooglePhotosProcessor{Source: source, Config: config}
	}
	log.Printf("Unrecognized source kind: %s\n", source.Kind)
	return nil
}

//Parses a line from the index and uses it to initialize a new MosaicTile
func createNodeFromLine(line string) (*gomosaic.MosaicTile, error) {
	// construct node
	parts := strings.Split(line, delimiter)
	if len(parts) != 5 {
		return nil, errors.New("invalid index line")
	}
	return &gomosaic.MosaicTile{Loc: parts[0], Filename: parts[1], AvgR: util.GetInt32(parts[2]),
		AvgG: util.GetInt32(parts[3]), AvgB: util.GetInt32(parts[4])}, nil

}
