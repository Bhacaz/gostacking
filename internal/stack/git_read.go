package stack

import (
	"errors"
	"github.com/Bhacaz/gostacking/internal/color"
)

func (sm StacksManager) currentBranchName() (string, error) {
	currentBranch, err := sm.gitExecutor.Exec("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", errors.New("failed to get current branch")
	}
	return currentBranch, nil
}

func (sm StacksManager) branchExists(branchName string) bool {
	_, err := sm.gitExecutor.Exec("rev-parse", "--verify", branchName)
	return err == nil
}

func (sm StacksManager) branchHasDiff(parentBranch string, branch string) (bool, error) {
	output, err := sm.gitExecutor.Exec("diff", "--name-only", branch+"..."+parentBranch)
	if err != nil {
		return false, errors.New("failed to get diff\n" + output)
	}
	return len(output) > 0, nil
}

func (sm StacksManager) lastLog(branch string) string {
	output, err := sm.gitExecutor.Exec("log", "--pretty=format:%s - %Cred%h%Creset - %C(bold blue)%an%Creset - %Cgreen%cr%Creset", "-n", "1", branch)
	if err != nil {
		return "could not get log"
	}
	return output
}

func (sm StacksManager) isBehindRemote(branch string) bool {
	output, err := sm.gitExecutor.Exec("diff", "--name-only", branch+"...origin/"+branch)
	if err != nil {
		return false
	}
	return len(output) > 0
}

func (sm StacksManager) fetch() error {
	_, err := sm.gitExecutor.Exec("fetch")
	if err != nil {
		return errors.New("failed to fetch")
	}
	return nil
}

func (sm StacksManager) gitClean() bool {
	output, err := sm.gitExecutor.Exec("status", "--porcelain")
	if err != nil {
		return false
	}
	return len(output) == 0
}

func (sm StacksManager) mainBranchWithRemote() (string, error) {
	main, err := sm.gitExecutor.Exec("symbolic-ref", "refs/remotes/origin/HEAD", "--short")

	if err != nil {
		return "", errors.New("Error getting origin default main branch:\n To set it try: " + color.Teal("git remote set-head origin <<main branch>>"))
	}
	return main, nil
}
