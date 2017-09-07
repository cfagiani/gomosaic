package indexer

import (
	"testing"
	"github.com/cfagiani/gomosaic/util"
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
