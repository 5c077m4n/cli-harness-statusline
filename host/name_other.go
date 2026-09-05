//go:build !darwin && !linux

package host

func processName(pid int) string {
	return ""
}
