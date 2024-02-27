package git

import (
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"log"
	"os/exec"
	"strings"
)

type interfaceGitExecutor interface {
	execCommand(command []string) (string, error)
}

type executor struct{}

func (e executor) execCommand(gitCmdArgs []string) (string, error) {
	execCmd := exec.Command("git", gitCmdArgs...)
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
	BranchDiff(baseBranch string, branch string) bool
	LastLog(branch string) string
}

type Commands struct {
	executor interfaceGitExecutor
}

func (c Commands) exec(command []string) (string, error) {
	return c.executor.execCommand(command)
}

func Cmd() Commands {
	return Commands{
		executor: executor{},
	}
}

func (c Commands) CurrentBranchName() (string, error) {
	currentBranch, err := c.exec([]string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return "", err
	}
	return currentBranch, nil
}

func (c Commands) BranchExists(branchName string) bool {
	_, err := c.exec([]string{"rev-parse", "--verify", "--quiet", "refs/heads/" + branchName})
	return err == nil
}

func (c Commands) Checkout(branchName string) {
	_, err := c.exec([]string{"checkout", branchName})
	if err != nil {
		log.Fatalf("Error checkout branch %s: %s", color.Yellow(branchName), err.Error())
	}
}

func (c Commands) SyncBranches(branches []string, checkoutBranchEnd string, push bool) {
	// Return if contains unstaged changes
	if !c.gitClean() {
		fmt.Println("Unstaged changes. Please commit or stash them.")
		return
	}

	fmt.Println("Fetching...")
	_, err := c.exec([]string{"fetch"})
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, branch := range branches {
		fmt.Println("Checkout to", color.Yellow(branch))
		_, err = c.exec([]string{"checkout", branch})
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Pulling", color.Yellow(branch), "...")
		_, err = c.exec([]string{"pull"})
		if err != nil {
			fmt.Println("Nothing to pull on", color.Yellow(branch))
		}

		// Nothing to merge on first branch
		if i == 0 {
			if push {
				c.pushBranch(branch)
			}
			continue
		}

		toMerge := branches[i-1]
		fmt.Println("Merging", color.Yellow(toMerge), "->", color.Yellow(branch))
		err = c.merge(branch, toMerge)
		if err != nil {
			fmt.Println(err)
			break
		}
		if push {
			c.pushBranch(branch)
		}
	}
	c.Checkout(checkoutBranchEnd)
}

func (c Commands) BranchDiff(baseBranch string, branch string) bool {
	output, err := c.exec([]string{"diff", "--name-only", branch + "..." + baseBranch})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return len(output) > 0
}

func (c Commands) LastLog(branch string) string {
	output, err := c.exec([]string{"log", "--pretty=format:%s - %Cred%h%Creset - %C(bold blue)%an%Creset - %Cgreen%cr%Creset", "-n", "1", branch})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return output
}

func (c Commands) pushBranch(branchName string) {
	fmt.Println("Pushing", color.Yellow(branchName), "...")
	_, err := c.exec([]string{"push"})
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commands) gitClean() bool {
	output, err := c.exec([]string{"status", "--porcelain"})
	if err != nil {
		return false
	}
	return len(output) == 0
}

func (c Commands) merge(currentBranch string, toMerge string) error {
	output, err := c.exec([]string{
		"merge",
		toMerge,
		"--no-squash",
		"--commit",
		"-m",
		"Merge branch " + toMerge + " into " + currentBranch + " (gostacking)",
	})
	if err != nil {
		fmt.Println("Error merging:", output)
		return err
	}
	return nil
}
