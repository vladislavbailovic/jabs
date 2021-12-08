package sys

import (
	"os"
	"testing"
)

func Test_TempFilePathReturnsUniquePath(t *testing.T) {
	temp1 := TempFilePath()
	stat1, err := os.Stat(temp1)
	if err != nil {
		t.Fatalf("Unable to stat temp file %s (1)", temp1)
	}
	mode1 := stat1.Mode()
	if !mode1.IsRegular() {
		t.Fatalf("Expected temp file to be a regular file: %s (1)", temp1)
	}

	temp2 := TempFilePath()
	stat2, err := os.Stat(temp1)
	if err != nil {
		t.Fatalf("Unable to stat temp file %s (2)", temp2)
	}
	mode2 := stat2.Mode()
	if !mode2.IsRegular() {
		t.Fatalf("Expected temp file to be a regular file: %s (2)", temp2)
	}

	if temp1 == temp2 {
		t.Fatalf("expected two temp files to differ: %s", temp1)
	}

	os.Remove(temp1)
	os.Remove(temp2)
}

func Test_NewFileishLoadsGlobFiles(t *testing.T) {
	expected := 2
	file := NewFileish("../examples/*.yml")
	if len(file.sources) != expected {
		t.Log(file.sources)
		t.Fatalf("expected exactly %d example yamls, got %d", expected, len(file.sources))
	}
}

func Test_NewFileishGlobRequiresYamlFiles(t *testing.T) {
	expected := 0
	file := NewFileish("../sys/*")
	if len(file.sources) != expected {
		t.Log(file.sources)
		t.Fatalf("expected exactly %d sys yamls, got %d", expected, len(file.sources))
	}
}

func Test_NewFileishLoadsSingleFile(t *testing.T) {
	expected := 1
	file := NewFileish("../examples/self.yml")
	if len(file.sources) != expected {
		t.Log(file.sources)
		t.Fatalf("expected exactly %d yamls, got %d", expected, len(file.sources))
	}
}

func Test_FileishReadsSingleYamlFile(t *testing.T) {
	path := "../examples/self.yml"
	stat1, _ := os.Stat(path)
	file := NewFileish(path)
	content := file.Read()
	if int64(len(content)) != stat1.Size() {
		t.Fatalf("expected %d bytes content, but got %d",
			stat1.Size(), len(content))
	}
}

func Test_FileishReadsGlobbedYamlFile(t *testing.T) {
	paths := []string{"../examples/self.yml", "../examples/first.yml"}
	var sumBytes int64
	for _, path := range paths {
		stat, _ := os.Stat(path)
		sumBytes += stat.Size()
	}
	file := NewFileish("../examples/*.yml")
	content := file.Read()
	if int64(len(content)) != sumBytes {
		t.Fatalf("expected %d bytes content, but got %d",
			sumBytes, len(content))
	}
}
