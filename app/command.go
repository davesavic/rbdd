package app

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func (a *APITest) iExecuteCommand(command string) error {
	return a.iExecuteCommandInDirectory(command, "")
}

func (a *APITest) iExecuteCommandInDirectory(command string, dir string) error {
	command = a.replaceVars(command)
	dir = a.replaceVars(dir)

	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("command is empty")
	}

	if a.debug {
		fmt.Printf("Executing command: %s\n", command)
		if dir != "" {
			fmt.Printf("In directory: %s\n", dir)
		}
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	a.commandOutput = strings.Trim(stdout.String(), "\n")
	if err != nil {
		return fmt.Errorf("command failed: %v\nStdout: %s\nStderr: %s",
			err, a.commandOutput, stderr.String())
	}

	if a.debug {
		fmt.Printf("Command output: %s\n", a.commandOutput)
	}

	return nil
}

func (a *APITest) iExecuteCommandWithTimeout(command string, timeoutSec int) error {
	done := make(chan error)

	go func() {
		done <- a.iExecuteCommand(command)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		return fmt.Errorf("command timed out after %d seconds: %s", timeoutSec, command)
	}
}

func (a *APITest) theCommandOutputShouldMatch(expected string) error {
	expected = a.replaceVars(expected)
	if a.commandOutput != expected {
		return fmt.Errorf("expected command output to be '%s', but got '%s'", expected, a.commandOutput)
	}
	return nil
}

func (a *APITest) theCommandOutputShouldContain(expected string) error {
	expected = a.replaceVars(expected)
	if !strings.Contains(a.commandOutput, expected) {
		return fmt.Errorf("expected command output to contain '%s', but got '%s'", expected, a.commandOutput)
	}
	return nil
}
