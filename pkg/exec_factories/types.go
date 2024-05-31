package exec_factories

import (
	"context"
	"os/exec"
	"testing"
)

type MockedCommandInvocationFn func(t *testing.T)

type CmdFactory = func(name string, arg ...string) *exec.Cmd

type CtxCmdFactory = func(ctx context.Context) CmdFactory
