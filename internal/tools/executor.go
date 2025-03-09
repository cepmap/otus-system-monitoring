package tools

import (
	"os/exec"
	"strings"
)

var ExecCommand = Exec

func Exec(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
