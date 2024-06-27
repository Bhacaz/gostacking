package cliexec

import (
	"github.com/Bhacaz/gostacking/internal/printer"
	"os/exec"
	"strings"
)

type InterfaceCliExecutor interface {
	Exec(command ...string) (string, error)
}

type Executor struct {
	baseCliCmd string
	verbose bool
	printer printer.Printer
}

func NewExecutor(baseCliCmd string, verbose bool) Executor {
	return Executor{
		baseCliCmd: baseCliCmd,
		verbose: verbose,
		printer: printer.NewPrinter(),
	}
}

func (e Executor) println(a ...interface{}) {
	if e.verbose {
		e.printer.Println(a...)
	}
}

func (e Executor) Exec(gitCmdArgs ...string) (string, error) {
	e.println("CMD:\t", "git", strings.Join(gitCmdArgs, " "))

	execCmd := exec.Command("git", gitCmdArgs...)
	output, err := execCmd.CombinedOutput()
	result := strings.TrimSuffix(string(output), "\n")

	e.println("OUTPUT:\t", result)
	if err != nil {
		e.println("ERROR:\t", err, "\n")
		return result, err
	}
	e.println("")

	return result, nil
}
