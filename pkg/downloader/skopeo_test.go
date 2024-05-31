// Skopeo copies an image to a tarball
package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var ctx context.Context = context.Background()

func TestSkopeoReturnErrCodeZero(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	os.Exit(0)
}

func TestSkopeoReturnErrCodeOne(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	os.Exit(1)
}

func TestSkopeoCommandLine(t *testing.T) {
	params := InitParams{
		CmdFactory: FakeExecCommand(TestSkopeoReturnErrCodeZero),
		BaseDir:    "/tmp/image",
	}
	copier := New(params)
	copier.SetTargetImage("registry/repo/name:tag")
	gotcmd, gotargs := copier.CommandLine()

	wantcmd := "skopeo"
	wantargs := []string{
		"copy",
		"--remove-signatures",
		"docker://registry/repo/name:tag",
		"docker-archive:///tmp/image/image.tar",
	}

	if gotcmd != wantcmd {
		t.Errorf("got %v, want %v", gotcmd, wantcmd)
	}

	if !reflect.DeepEqual(gotargs, wantargs) {
		t.Errorf("got %v, want %v", gotargs, wantargs)
	}
}

func TestSkopeoCommandRunSuccess(t *testing.T) {
	params := InitParams{
		CmdFactory: FakeExecCommand(TestSkopeoReturnErrCodeZero),
		BaseDir:    "/tmp/image",
	}
	copier := New(params)

	copier.SetTargetImage("registry/repo/name:tag")

	if err := copier.Run(); err != nil {
		t.Errorf("command failed even if we marked as successful??")
	}
}

func TestSkopeoCommandRunFail(t *testing.T) {
	params := InitParams{
		CmdFactory: FakeExecCommand(TestSkopeoReturnErrCodeOne),
		BaseDir:    "/tmp/image",
	}
	copier := New(params)
	copier.SetTargetImage("registry/repo/name:tag")

	if err := copier.Run(); err == nil {
		t.Errorf("command worked even if we marked as failing??")
	}
}

func TestSkopeoCommandWithCredentials(t *testing.T) {
	params := InitParams{
		CmdFactory: FakeExecCommand(TestSkopeoReturnErrCodeZero),
		BaseDir:    "/tmp/image",
	}
	copier := New(params)

	os.Setenv("REGISTRY_USERNAME", "user")
	os.Setenv("REGISTRY_PASSWORD", "password")
	defer os.Unsetenv("REGISTRY_USERNAME")
	defer os.Unsetenv("REGISTRY_PASSWORD")

	copier.SetTargetImage("registry/repo/name:tag")
	gotcmd, gotargs := copier.CommandLine()

	wantcmd := "skopeo"
	wantargs := []string{
		"copy",
		"--remove-signatures",
		"--src-username",
		"user",
		"--src-password",
		"password",
		"docker://registry/repo/name:tag",
		"docker-archive:///tmp/image/image.tar",
	}

	if gotcmd != wantcmd {
		t.Errorf("got %v, want %v", gotcmd, wantcmd)
	}

	if !reflect.DeepEqual(gotargs, wantargs) {
		t.Errorf("got %v, want %v", gotargs, wantargs)
	}
}
