package sys

import "testing"

func Test_Shell(t *testing.T) {
	res := Shell()
	if res != "/bin/bash" {
		t.Fatalf("expected bash res")
	}
}

func Test_ShellCommandParam(t *testing.T) {
	res := ShellCommandParam()
	if res != "-c" {
		t.Fatalf("expected -c")
	}
}

func Test_Shebang(t *testing.T) {
	res := Shebang()
	if res != "#!/bin/bash" {
		t.Fatalf("expected bash shebang")
	}
}
