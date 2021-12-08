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
