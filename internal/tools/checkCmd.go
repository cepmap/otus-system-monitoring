package tools

import (
	"fmt"
	"os/exec"
)

// CheckCommand проверяет наличие команды в системе
func CheckCommand(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("command %s not found: %w", name, err)
	}
	return nil
}
