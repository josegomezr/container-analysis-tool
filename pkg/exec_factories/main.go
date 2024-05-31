package exec_factories

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

var DefaultCmdFactory CmdFactory = exec.Command

var ContextualCmdFactory CtxCmdFactory = func(ctx context.Context) CmdFactory {
	return func(name string, arg ...string) *exec.Cmd {
		return exec.CommandContext(ctx, name, arg...)
	}
}

func FakeExecCommand(fn MockedCommandInvocationFn) CmdFactory {
	components := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	name := components[len(components)-1]

	return func(command string, args ...string) *exec.Cmd {
		cs := []string{fmt.Sprintf("-test.run=%s", name), "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}
