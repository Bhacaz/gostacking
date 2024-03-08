package git

import (
	"github.com/Bhacaz/gostacking/internal/printer"
	"os/exec"
	"strings"
)

type InterfaceGitExecutor interface {
	Exec(command ...string) (string, error)
}

type Executor struct {
	verbose bool
	printer printer.Printer
}

func NewExecutor(verbose bool) Executor {
	return Executor{
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
