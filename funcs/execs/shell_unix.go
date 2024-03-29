//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || plan9 || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd plan9 solaris

package execs

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

type shellResult struct {
	output string
	err    error
}

// ExecShellWithContext Execute the shell command to set the execution timeout
func ExecShellWithContext(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	resultChan := make(chan shellResult)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				resultChan <- shellResult{"", errors.New(fmt.Sprintf("%v", err))}
			}
		}()
		output, err := cmd.CombinedOutput()
		resultChan <- shellResult{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return result.output, result.err
	}
}
