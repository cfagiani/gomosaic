package util

import (
	"testing"
	"os"
)

//TestReadConfig validates that the ReadConfig function correctly populates a Config struct when given a valid
//configuration file.
func TestReadConfig(t *testing.T) {
	cases := []struct {
		source      string
		expectError bool
	}{
		{"../testdata/testconfig.json", false},
		{"../testdata/testtoken.json", true},
		{"../testdata/missingfile", true},
	}
	for _, c := range cases {
		config, err := ReadConfig(c.source)
		if err != nil && !c.expectError {
			t.Errorf("ReadConfig returned an unexpected error for %v", c.source)
		} else if err == nil && c.expectError {
			t.Errorf("Expected ReadConfig to return an error for %v but it did not", c.source)
		} else if err == nil {
			if config.GoogleClientId != "thisisnotreal" || len(config.Sources) != 1 {
				t.Errorf("ReadConfig did not populate the config object correctly for %v", c.source)
			}
		}
	}
}

//TestGetPath validates the behavior of GetPath by ensuring the directory is combined with the filename using the os
//pathSeparator
func TestGetPath(t *testing.T) {
	cases := []struct {
		dir      string
		file     string
		expected string
	}{
		{"", "test", string(os.PathSeparator) + "test"},
		{"blah", "test", "blah" + string(os.PathSeparator) + "test"},
		{"blah", "", "blah" + string(os.PathSeparator)},
	}
	for _, c := range cases {
		path := GetPath(c.dir, c.file)
		if path != c.expected {
			t.Errorf("GetPath returned %v but should have returned %v", path, c.expected)
		}
	}
}

//TestGetInt32 verifies that GetInt32 eats errors when converting from strings containing numbers to unsigned ints.
func TestGetInt32(t *testing.T) {
	cases := []struct {
		input    string
		expected uint32
	}{
		{"123", 123},
		{"0123", 123},
		{"-123", 0},
		{"NaN", 0},
	}
	for _, c := range cases {
		val := GetInt32(c.input)
		if val != c.expected {
			t.Errorf("GetInt32 returned %v but should have returned %v", val, c.expected)
		}
	}
}

func TestGetPhotosService(t *testing.T) {

	client, err := GetPhotosService("a", "b", "../testdata/testtoken.json")
	if err != nil {
		t.Errorf("GetPhotosService returned an unexpected error %v", err)
	} else {
		if client == nil {
			t.Error("GetPhotosService returned nil unexpectedly")
		}
	}
}
