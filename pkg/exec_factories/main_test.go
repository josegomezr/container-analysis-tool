package exec_factories

import (
	"fmt"
	"os"
	"testing"
)

const exampleOutput = "result"

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, exampleOutput)
	os.Exit(0)
}

func TestRunCommand(t *testing.T) {
	old := DefaultCmdFactory
	DefaultCmdFactory = FakeExecCommand(TestHelperProcess)
	defer func() { DefaultCmdFactory = old }()

	out, err := DefaultCmdFactory("my command example").Output()

	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}

	if string(out) != exampleOutput {
		t.Errorf("Expected %q, got %q", exampleOutput, out)
	}
}
