// +build windows

package exec

import (
	"fmt"
	"os/exec"
	"strconv"
)

func MakeCmd(args ...string) *exec.Cmd {
	return exec.Command("cmd", append([]string{"/C"}, args...)...)
}

func KillProcess(cmd *exec.Cmd) error {
	return exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
}

func OpenInBrowser(path string) error {
	logs, err := MakeCmd("start", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w, logs: %q", err, logs)
	}
	return nil
}
