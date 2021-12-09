package cmd

import (
	"testing"
	"jabs/types"
)

func Test_GetDirs(t *testing.T) {
	dirs := getDirs()
	expected := 8
	if len(dirs) != expected {
		t.Log(dirs)
		t.Fatalf("expected exactly %d dirs, got %d", expected, len(dirs))
	}
}

func Test_FilenamePatternFilter(t *testing.T) {
	pf := FilenameFilter{types.FilenamePattern("*.go")}

	if pf.Matches("whatever/should/not.match") {
		t.Fatalf("expected random crap to not match path pattern")
	}

	if pf.Matches("relative/path/to/file.go") {
		t.Fatalf("relative path expected to NOT match, but it did")
	}

	if !pf.Matches("file.go") {
		t.Fatalf("just filename expected to match, but it didn't")
	}
}
