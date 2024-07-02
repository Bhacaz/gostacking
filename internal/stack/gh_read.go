package stack

import (
	"errors"
	"strings"
)

func (sm StacksManager) ghCliConfigure() error {
	output, err := sm.ghExecutor.Exec("auth", "status")
	if err != nil && strings.Contains(output, "not found") {
		return errors.New("GH-CLI not found")
	}
	if strings.Contains(output, "You are not logged") {
		return errors.New(output)
	}
    if err != nil {
        return err
    }
	return nil
}

func (sm StacksManager) ghPrNumber(branchName string) (string, error) {
	output, err := sm.ghExecutor.Exec("pr", "view", branchName, "-q", ".number", "--json=number")

	if err != nil {
		return "", err
	}

	return output, nil
}
