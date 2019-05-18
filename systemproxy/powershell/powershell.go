package powershell

import (
	"bytes"
	"os/exec"
)

type PowerShell struct {
	powerShell string
}

func New() *PowerShell {
	ps, err := exec.LookPath("powershell.exe")

	if err != nil {
		println("powershell not found")
		return nil
	}

	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) Execute(arg string) (stdOut string, stdErr string, err error) {
	//fmt.Printf("executing\n %v\n", arg)
	//args := append([]string{"-NoProfile", "-NonInteractive"}, arg)
	cmd := exec.Command(p.powerShell, arg)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}
