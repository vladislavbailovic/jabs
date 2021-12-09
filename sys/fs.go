package sys

import (
	"io/ioutil"
	"os"
	"jabs/types"
	"jabs/dbg"
	"path/filepath"
)

const (
	TEMP_DIR          string = "" // Use system temp dir
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
	path    string
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

func (f Fileish) Read() []byte {
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

type Dirlist struct {
	root types.PathPattern
}

func (d Dirlist) Read() []string {
	return d.readDirs(d.root, []string{})
}

func (d Dirlist) readDirs(root types.PathPattern, dirs []string) []string {
	dirFilter := DirFilter{}
	paths, err := filepath.Glob(string(root))
	if err != nil {
		dbg.FatalError("Unable to get directories (%s): %v", root, err)
	}

	for _, path := range paths {
		if filepath.Base(path)[0:1] == "." {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			dbg.FatalError("Unable to determine path for %s: %v", path, err)
		}

		if dirFilter.Matches(abs) {
			dirs = append(dirs, abs)
			dirs = d.readDirs(
				types.PathPattern(abs + string(os.PathSeparator) + "*"), dirs)
		}
	}

	return dirs
}

func NewDirlist(root types.PathPattern) Dirlist {
	return Dirlist{root}
}

type Filter interface {
	Matches(path string) bool
}

type FilenameFilter struct {
	pattern types.FilenamePattern
}

func (pf FilenameFilter)Matches (what string) bool {
	isMatch, err := filepath.Match(string(pf.pattern), what)
	if err != nil {
		return false
	}
	return isMatch
}

type DirFilter struct {}
func (df DirFilter)Matches (what string) bool {
	stat, err := os.Stat(what)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func NewFilter(params ...interface{}) Filter {
	if len(params) > 0 {
		return FilenameFilter{params[0].(types.FilenamePattern)}
	}
	return DirFilter{}
}
