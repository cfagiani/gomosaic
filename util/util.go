package util

import (
	"os"
	"log"
	"strconv"
)

func GetAbsolutePath(dir string, file string) string {
	return dir + string(os.PathSeparator) + file
}


func CheckError(err error, msg string, isFatal bool) bool {
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

// Converts a string to an integer, eating any errors
func GetInt(s string) uint32 {
	// TODO: actually handle the error
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint32(i)
	}
	return 0
}
