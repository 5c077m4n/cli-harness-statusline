// Package host detects the CLI that invoked the statusline.
package host

import "os"

// Name returns the process name of the parent process (the host CLI).
func Name() string {
	return processName(os.Getppid())
}
