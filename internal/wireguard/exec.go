package wireguard

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// execCommand executes a shell command and returns its output
func execCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("command failed: %s: %s", command, errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// execCommandWithInput executes a shell command with input piped to stdin
func execCommandWithInput(input, command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("command failed: %s: %s", command, errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}
