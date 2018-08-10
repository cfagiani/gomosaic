package indexer

import (
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/indexer/processor"
	"github.com/cfagiani/gomosaic/util"
	"reflect"
	"testing"
	"os"
	"strings"
)

func TestGetIndexFileName(t *testing.T) {
	cases := []struct {
		source       string
		expectedName string
		expectExist  bool
	}{
		{".", util.GetPath(".", "mosaicIndex.dat"), false},
		{"indexer_test.go", "indexer_test.go", true},
		{"/should/not/exist/or/it/will/fail", "/should/not/exist/or/it/will/fail", false},
	}
	for _, c := range cases {
		name, exist := GetIndexFileName(c.source)
		if name != c.expectedName || exist != c.expectExist {
			t.Errorf("GetIndexFilname(%q) == %q,%t want %q,%t",
				c.source, name, exist, c.expectedName, c.expectExist)
		}
	}
}

//TestGetProcessor ensures that the GetProcessor method returns the right type of IndexProcessor instance based on the specified type.
func TestGetProcessor(t *testing.T) {
	cases := []struct {
		kind         string
		expectedType processor.IndexProcessor
	}{
		{kind: processor.LocalKind, expectedType: processor.LocalProcessor{Source: gomosaic.ImageSource{}}},
		{kind: processor.GoogleKind, expectedType: processor.GooglePhotosProcessor{Config: gomosaic.Config{},
			Source: gomosaic.ImageSource{}}},
		{kind: "junk", expectedType: nil},
		{kind: "", expectedType: nil},
	}
	for _, c := range cases {
		processor := getProcessor(gomosaic.ImageSource{Kind: c.kind, Path: "", Options: ""}, gomosaic.Config{})
		if reflect.TypeOf(processor) != reflect.TypeOf(c.expectedType) {
			t.Errorf("GetProcessor for %q returned wrong type. Got %q wanted %q", c.kind,
				reflect.TypeOf(processor), reflect.TypeOf(c.expectedType))
		}
	}
}

//TestReadIndex verifies the ReadIndex function
func TestReadIndex(t *testing.T) {
	cases := []struct {
		filename string
		count    int
	}{
		{"../testdata/testindex.dat", 4},
		{"../testdata/notthere", 0},
		{"../testdata/img1.png", 0},
	}
	for _, c := range cases {
		index := ReadIndex(c.filename)
		if c.count != index.Len() {
			t.Errorf("ReadIndex returned an index with %d entries for %s. Wanted %d", index.Len(), c.filename, c.count)
		}
	}
}

//TestIndex verifies that we can index files and that when we re-read it, it contains what it should.
func TestIndex(t *testing.T) {
	destName := "../testdata/tempindex.dat"
	expectedCount := 4
	defer func() {
		//cleanup no matter what.
		os.Remove(destName)
	}()
	err := Index("../testdata/testconfig.json", destName)
	if err != nil {
		t.Errorf("Could not index files %v", err)
	}
	index := ReadIndex(destName)
	if index.Len() != expectedCount {
		t.Errorf("Expected to index %d files but found %d", expectedCount, index.Len())
	}
	var prevFile = ""
	for _, tile := range index {
		if tile.Loc != "L" {
			t.Errorf("Expected index to only contain local file but found %s", tile.Loc)
		}
		if prevFile != "" {
			if strings.Compare(prevFile, tile.Filename) > 0 {
				t.Errorf("Expected index to be sorted by filename but %s came before %s", prevFile, tile.Filename)
			}
		}
		prevFile = tile.Filename
	}
}
