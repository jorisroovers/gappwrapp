package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func GappWrapp(command []string) (int, string, string) {
	// prepend "--" to make it clear to gappwrapp that any string that follows is part of the command to execute
	command = append([]string{"--"}, command...)
	cmd := exec.Command("bin/darwin/gappwrapp", command...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	exitCode := 0
	exitError, ok := err.(*exec.ExitError)
	if ok {
		exitCode = exitError.ExitCode()
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	return exitCode, outStr, errStr
}

func TestBasics(t *testing.T) {
	exitCode, outStr, errStr := GappWrapp([]string{"./test-script.sh", "37", "foobar", "two words", "hurdur"})
	if exitCode != 37 {
		t.Errorf("Incorrect exit code, expected '%v', got '%v'", 37, exitCode)
	}
	expectedStdout := "Hello to stdout\nARG1 37\nARG2 foobar\nARG3 two words\nARG4 hurdur\n"
	if outStr != expectedStdout {
		t.Errorf("Incorrect stdout output, expected '%v', got '%v'", expectedStdout, outStr)
	}
	expectedStderr := "Hello to stderr\n"
	if errStr != expectedStderr {
		t.Errorf("Incorrect stderr output, expected '%v', got '%v'", expectedStderr, errStr)
	}
}
