package sys

import (
	"io/ioutil"
	"path/filepath"
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

type Fileish struct {
	path string
	sources []string
}

func NewFileish(path string) Fileish {
	matches, err := filepath.Glob(path)
	if err != nil {
		dbg.Error("Unable to resolve path spec: %s (%v)", path, err)
	}
	if len(matches) == 0 {
		dbg.Error("Could not resolve path spec to concrete files: %s", path)
	}
	sources := []string{}
	for _, match := range matches {
		abs, err := filepath.Abs(match)
		if err != nil {
			dbg.Warning("Could not resolve absolute path for %s: %v", match, err)
			continue
		}
		ext := filepath.Ext(abs)
		if ext != ".yml" && ext != ".yaml" {
			dbg.Warning("Not a yaml file by extension: %s", abs)
			continue
		}
		sources = append(sources, abs)
	}
	if len(sources) == 0 {
		dbg.Warning("Resolving sources came up empty")
	}
	return Fileish{path, sources}
}

func (f Fileish)Read() []byte {
	data := []byte{}
	for _, source := range f.sources {
		file, err := ioutil.ReadFile(source)
		if err != nil {
			dbg.Warning("Unable to read source: %s (%v)", source, err)
			continue
		}
		data = append(data, file...)
	}
	return data
}
