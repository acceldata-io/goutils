// Package cmd is simple wrapper on top of the in-built os/exec module
package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"
)

// Command is an instance of an executable command
type Command struct {
	Command string
	Args    []string
	Status  Status
	Ctx     context.Context
}

// Status contains information about an executed command instance
type Status struct {
	Process  *os.Process
	ExitCode int
	StdOut   string
	StdErr   string
}

// New returns a new command instance
func New(ctx context.Context, command string, args []string) *Command {
	return &Command{
		Command: command,
		Args:    args,
		Status:  Status{},
		Ctx:     ctx,
	}
}

// Run execute the command and returns the command execution status, stdout and stderr
func (c *Command) Run() (*Command, error) {
	cmd := exec.CommandContext(c.Ctx, c.Command, c.Args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return c, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return c, err
	}

	if err := cmd.Start(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			c.Status.ExitCode = exitError.ExitCode()
		}
	}

	stdoutByte, err := io.ReadAll(stdout)
	if err != nil {
		return c, err
	}

	stderrByte, err := io.ReadAll(stderr)
	if err != nil {
		return c, err
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			c.Status.ExitCode = exitError.ExitCode()
		}
	}
	c.Status.Process = cmd.Process
	c.Status.StdOut = string(stdoutByte)
	c.Status.StdErr = string(stderrByte)
	return c, nil
}

// NewWithExpression creates a new Command with the specified command binary and the expression
func (c *Command) NewWithExpression(cmdBin string, expression string) *Command {
	c.Command = cmdBin
	c.Args = []string{"-c", expression}

	return c
}
