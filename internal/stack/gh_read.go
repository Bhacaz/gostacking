package stack

import (
	"errors"
	"github.com/Bhacaz/gostacking/internal/cliexec"
	"strings"
)

func (sm StacksManager) ghCliConfigure() error {
	var ghCliExecutor = cliexec.NewExecutor("gh", false)
	output, err := ghCliExecutor.Exec("auth", "status")
	if err != nil && strings.Contains(output, "executable file not found") {
		return errors.New("GH-CLI not found")
	}
	if err != nil && strings.Contains(output, "You are not logged") {
		return errors.New(output)
	}
	return nil
}

func (sm StacksManager) ghPrNumber(branchName string) (string, error) {
	var ghCliExecutor = cliexec.NewExecutor("gh", false)
	output, err := ghCliExecutor.Exec("pr", "view", branchName, "-q", ".number", "--json=number")

	if err != nil {
		return "", err
	}

	return output, nil
}