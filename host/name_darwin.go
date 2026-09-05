//go:build darwin

package host

import "golang.org/x/sys/unix"

func processName(pid int) string {
	proc, err := unix.SysctlKinfoProc("kern.proc.pid", pid)
	if err != nil {
		return ""
	}
	comm := proc.Proc.P_comm
	var name []byte
	for _, b := range comm {
		if b == 0 {
			break
		}
		name = append(name, b)
	}
	return string(name)
}
