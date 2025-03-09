package tools

import (
	"fmt"
	"os/exec"
	"strconv"
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

func CheckCommand(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("command %s not found: %w", name, err)
	}
	return nil
}

func ParseFloat(input string) float64 {
	buff := strings.ReplaceAll(input, ",", ".")
	output, err := strconv.ParseFloat(buff, 64)
	if err != nil {
		return 0.0
	}
	return output
}
