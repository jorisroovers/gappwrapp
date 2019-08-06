package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	cli "github.com/jawher/mow.cli"
)

func main() {
	cp := cli.App("gappwrapp", "Lightweight process wrapper to execute commands with user-defined side-effects.")
	cp.Spec = "[--debug] CMD..." // '...' = multi-string arg
	cp.Version("v version", "gappwrapp 0.1")
	debug := cp.BoolOpt("debug", false, "Enable debug mode")
	// CMD = multi-string argument
	command := cp.StringsArg("CMD", nil, "Command to execute")

	cp.Action = func() {
		// setup logger if debug is enabled, disable otherwise
		if *debug {
			log.SetOutput(os.Stderr)
			log.Println("Debug mode enabled")
		} else { // as per https://stackoverflow.com/a/34457930/381010
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
		}

		log.Printf("Passed command: %#v", *command)
		exitCode := runCommand(command)

		log.Println("Command exit code:", exitCode)
		os.Exit(exitCode)
	}

	cp.Run(os.Args)
}

func readAndNotify(reader io.ReadCloser, writer io.Writer, waitgroup *sync.WaitGroup) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Fprintf(writer, "%s\n", scanner.Text())
	}
	waitgroup.Done()
}

func runCommand(command *[]string) int {

	cmdName := (*command)[0]
	cmdArgs := (*command)[1:]

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Env = os.Environ() // pass along the current environment
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go readAndNotify(stdoutReader, os.Stdout, &wg)
	go readAndNotify(stderrReader, os.Stderr, &wg)

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	}

	err = cmd.Wait()
	// Get exit code. Requires some minor vodoo, although improved since Go 1.12
	// https://stackoverflow.com/a/55055100/381010
	exitCode := 0
	exitError, ok := err.(*exec.ExitError)
	if ok {
		exitCode = exitError.ExitCode()
	}

	log.Println("Waiting for stdout/stderr output to finish")
	wg.Wait()
	return exitCode
}
