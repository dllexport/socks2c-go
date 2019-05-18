package systemproxy

import "os/exec"

func execAndGetRes(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return string(out[:])
}
