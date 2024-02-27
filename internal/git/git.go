package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type interfaceGitExecutor interface {
	execCommand(command string) (string, error)
}

type executor struct{}

func (e executor) execCommand(gitCmd string) (string, error) {
	cmdArgs := strings.Fields(gitCmd)
	execCmd := exec.Command("git", cmdArgs...)
	output, err := execCmd.CombinedOutput()
	result := strings.TrimSuffix(string(output), "\n")

	if err != nil {
		//fmt.Println("Git gitCmd err:", result)
		return "", err
	}

	return result, nil
}

type InterfaceCommands interface {
	CurrentBranchName() (string, error)
	BranchExists(branchName string) bool
	Checkout(branchName string)
	SyncBranches(branches []string, checkoutBranchEnd string, push bool)
}

type Commands struct {
	executor interfaceGitExecutor
}

func (c Commands) exec(command string) (string, error) {
	return c.executor.execCommand(command)
}

func Cmd() Commands {
	return Commands{
		executor: executor{},
	}
}

func (c Commands) CurrentBranchName() (string, error) {
	currentBranch, err := c.exec("rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", err
	}
	return currentBranch, nil
}

func (c Commands) BranchExists(branchName string) bool {
	_, err := c.exec("rev-parse --verify --quiet \"refs/heads/" + branchName + "\"")
	return err == nil
}

func (c Commands) Checkout(branchName string) {
	_, err := c.exec("checkout " + branchName)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commands) SyncBranches(branches []string, checkoutBranchEnd string, push bool) {
	// Return if contains unstaged changes
	if !c.gitClean() {
		fmt.Println("Unstaged changes. Please commit or stash them.")
		return
	}

	fmt.Println("Fetching...")
	_, err := c.exec("fetch")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, branch := range branches {
		fmt.Println("Checkout to", branch)
		_, err := c.exec("checkout " + branch)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Pulling", branch, "...")
		_, err = c.exec("pull")
		if err != nil {
			fmt.Println(err)
		}

		// Nothing to merge on first branch
		if i == 0 {
			if push {
				c.pushBranch(branch)
			}
			continue
		}

		toMerge := branches[i-1]
		fmt.Println("Merging", toMerge, "->", branch)
		err = c.merge(branch, toMerge)
		if err != nil {
			fmt.Println(err)
			break
		}
		if push {
			c.pushBranch(branch)
		}
	}
	_, err = c.exec("checkout " + checkoutBranchEnd)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commands) pushBranch(branchName string) {
	fmt.Println("Pushing", branchName, "...")
	_, err := c.exec("push")
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commands) gitClean() bool {
	output, err := c.exec("status --porcelain")
	if err != nil {
		return false
	}
	return len(output) == 0
}

func (c Commands) merge(currentBranch string, toMerge string) error {
	output, err := c.exec("merge " + toMerge + " --no-squash --commit -m \"Merge branch " + toMerge + " into " + currentBranch + " (gostacking)\"")
	if err != nil {
		fmt.Println("Error merging:", output)
		return err
	}
	return nil
}
