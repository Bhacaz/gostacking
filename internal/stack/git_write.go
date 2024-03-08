package stack

import (
	"errors"
	"github.com/Bhacaz/gostacking/internal/color"
	"strings"
)

func (sm StacksManager) checkout(branchName string) error {
	output, err := sm.gitExecutor.Exec("checkout", branchName)
	if err != nil {
		return errors.New("failed to checkout " + color.Yellow(branchName) + "\n" + output)
	}
	return nil
}

// pushBranch will push the current branch to the remote
// If the branch does not have a remote, it will NOT return an error
func (sm StacksManager) pushBranch() error {
	output, err := sm.gitExecutor.Exec("push")
	if err != nil && !strings.Contains(output, "has no upstream branch") {
		return errors.New("failed to push\n" + output)
	}
	return nil
}

// pullBranch will pull the current branch from the remote
// If the branch does not have a remote, it will NOT return an error
func (sm StacksManager) pullBranch() error {
	output, err := sm.gitExecutor.Exec("pull")
	if err != nil && !strings.Contains(output, "There is no tracking information") {
		return errors.New("failed to pull\n" + output)
	}
	return nil
}

func (sm StacksManager) merge(currentBranch string, parentBranch string) error {
	output, err := sm.gitExecutor.Exec(
		"merge",
		parentBranch,
		"-m",
		"Merge branch "+parentBranch+" into "+currentBranch+" (gostacking)",
	)
	if err != nil {
		return errors.New("failed to merge\n" + output)
	}
	return nil
}
