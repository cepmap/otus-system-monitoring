package float

import (
	"errors"
	"os/exec"
	"testing"
)

type TestCommand struct {
	Output string
	Err    error
}

func (m *TestCommand) Run() (string, error) {
	return m.Output, m.Err
}

type TestCommandCreator func(command string, args ...string) *TestCommand

var commandCreator TestCommandCreator

func init() {
	commandCreator = func(command string, args ...string) *TestCommand {
		cmd := exec.Command(command, args...)
		output, err := cmd.CombinedOutput()
		return &TestCommand{
			Output: string(output),
			Err:    err,
		}
	}
}

func ExecTest(command string, args string) (string, error) {
	cmd := commandCreator(command, args)
	return cmd.Run()
}

func TestExec(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        string
		output      string
		error       error
		expected    string
		expectError bool
	}{
		{
			name:        "successful execution",
			command:     "echo",
			args:        "hello",
			output:      "hello",
			error:       nil,
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "command failure",
			command:     "false",
			args:        "",
			output:      "",
			error:       errors.New("exit status 1"),
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid command",
			command:     "nonexistent",
			args:        "",
			output:      "",
			error:       errors.New("executable file not found in $PATH"),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandCreator = func(command string, args ...string) *TestCommand {
				return &TestCommand{
					Output: tt.output,
					Err:    tt.error,
				}
			}

			output, err := ExecTest(tt.command, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if output != tt.expected {
				t.Errorf("expected output: %q, got: %q", tt.expected, output)
			}
		})
	}
}
