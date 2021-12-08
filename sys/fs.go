package sys

import (
	"io/ioutil"
	"jabs/dbg"
)

const (
	TEMP_DIR string = "" // Use system temp dir
	TEMP_FILE_PATTERN string = "jabs-tmp-*"
)

func TempFilePath() string {
	file, err := ioutil.TempFile(TEMP_DIR, TEMP_FILE_PATTERN)
	if err != nil {
		dbg.Error("Unable to create temp file: %v", err)
	}
	file.Close()
	return file.Name()
}
