//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !plan9 && !solaris
// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!plan9,!solaris

package execs

import (
	"context"
	"errors"
)

// Execute the shell command to set the execution timeout
func ExecShellWithContext(ctx context.Context, command string) (string, error) {
	return "", errors.New("this function does not support running under the current system")
}
