package execs

import (
	"bytes"
	"os/exec"
)

// Blocking functions that execute external shell commands,
// wait for the execution to be completed and return standard output
func Shell(name string, arg ...string) ([]byte, error) {
	// The function returns a *cmd, which is used to execute the program specified by name with the given parameters
	cmd := exec.Command(name, arg...)

	// Read cmd.stdout of io.writer type, and then convert byte type to [] byte type through bytes.buffer (buffer of byte type)
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run executes the commands contained in C and blocks until it is completed. Here stdout is taken out,
	// and cmd.wait () cannot get stdin, stdout, stderr correctly, so it is blocked there.
	if err := cmd.Run(); err != nil {
		return nil, err
	} else {
		return out.Bytes(), nil
	}
}
