package rund

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type (
	// Operator defines the Run interface to run operation.
	Operator interface {
		Run() error
	}

	// NoopOperator does nothing.
	NoopOperator struct{}

	// CmdOperator runs external command.
	CmdOperator struct {
		envar map[string]string
		cmd   []string
	}

	// FuncOperator executes function.
	FuncOperator func() error
)

var errMissingCmd = errors.New("missing Cmd")

// Run runs nooperation.
func (op NoopOperator) Run() error { return nil }

// Run runs external command.
func (op CmdOperator) Run() error {
	cmd := exec.Command(op.cmd[0], op.cmd[1:]...)
	env := os.Environ()

	for k, v := range op.envar {
		env = append(env,
			fmt.Sprintf("%s=%s", strings.ToUpper(k), v))
	}

	cmd.Env = env
	return cmd.Run()
}

// Run executes function.
func (op FuncOperator) Run() error { return op() }

// NewCmdOperator returns a new CmdOperator.
func NewCmdOperator(cmd []string, envar map[string]string) (CmdOperator, error) {
	if len(cmd) < 1 {
		return CmdOperator{}, errMissingCmd
	}
	return CmdOperator{
		cmd:   cmd,
		envar: envar,
	}, nil
}

// NewNoopOperator returns a new NoopOperator.
func NewNoopOperator() NoopOperator {
	return NoopOperator{}
}

// NewFuncOperator returns a new FuncOperator.
func NewFuncOperator(f func() error) FuncOperator {
	return FuncOperator(f)
}
