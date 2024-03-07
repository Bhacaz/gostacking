package stack

import (
	"errors"
	"github.com/Bhacaz/gostacking/internal/color"
)

func (sm StacksManager) checkout(branchName string) error {
	_, err := sm.gitExecutor.Exec("checkout", branchName)
	if err == nil {
		return errors.New("failed to checkout " + color.Yellow(branchName))
	}
	return nil
}

func (sm StacksManager) pushBranch() {
	_, _ = sm.gitExecutor.Exec("push")
}

func (sm StacksManager) merge(currentBranch string, parentBranch string) error {
	output, err := sm.gitExecutor.Exec(
		"merge",
		parentBranch,
		"--no-squash",
		"--commit",
		"-m",
		"Merge branch "+parentBranch+" into "+currentBranch+" (gostacking)",
	)
	if err != nil {
		return errors.New("failed to merge\n" + output)
	}
	return nil
}
