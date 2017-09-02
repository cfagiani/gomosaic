package util

import (
	"os"
	"log"
	"strconv"
)

//utility to get the absolute path of a file by concatenating the directory, the pathSeparator, and the file name.
func GetAbsolutePath(dir string, file string) string {
	return dir + string(os.PathSeparator) + file
}

//Checks the error argument and, if it is not nil, it will log the msg passed in. If isFatal is true, the log will be
//written as Fatal which will cause exit(1) to be called.
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

// Converts a string to a 32-bit unsigned integer, eating any errors
func GetInt(s string) uint32 {
	// TODO: actually handle the error?
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return uint32(i)
	}
	return 0
}
