package constants

import "os"

var INDEX_FILE string
var STORAGE_LOC string

func LoadConstants() error {
	INDEX_FILE = os.Getenv("INDEX_FILE")
	STORAGE_LOC = os.Getenv("STORAGE_LOC")
	return nil
}
