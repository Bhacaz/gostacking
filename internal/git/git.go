package git

import (
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"log"
	"os/exec"
	"strings"
)

type interfaceGitExecutor interface {
	execCommand(command ...string) (string, error)
}

type executor struct{}

func (e executor) execCommand(gitCmdArgs ...string) (string, error) {
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
	IsBehindRemote(branch string) bool
	Fetch()
}

type Commands struct {
	executor interfaceGitExecutor
}

func (c Commands) exec(command ...string) (string, error) {
	return c.executor.execCommand(command...)
}

func Cmd() Commands {
	return Commands{
		executor: executor{},
	}
}

func (c Commands) CurrentBranchName() (string, error) {
	currentBranch, err := c.exec("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return currentBranch, nil
}

func (c Commands) BranchExists(branchName string) bool {
	_, err := c.exec("rev-parse", "--verify", "--quiet", "refs/heads/"+branchName)
	return err == nil
}

func (c Commands) Checkout(branchName string) {
	_, err := c.exec("checkout", branchName)
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

	fmt.Println("Fetching", "...")
	c.Fetch()

	for i, branch := range branches {
		fmt.Println("Branch:", color.Yellow(branch))
		fmt.Println("\tCheckout", "...")
		_, err := c.exec("checkout", branch)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("\tPull", "...")
		_, err = c.exec("pull")
		if err != nil {
			fmt.Println("\tNothing to pull on", color.Yellow(branch))
		}

		// Nothing to merge on first branch
		if i == 0 {
			if push {
				c.pushBranch(branch)
			}
			continue
		}

		toMerge := branches[i-1]
		fmt.Println("\tMerging", color.Yellow(toMerge))
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
	output, err := c.exec("diff", "--name-only", branch+"..."+baseBranch)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return len(output) > 0
}

func (c Commands) LastLog(branch string) string {
	output, err := c.exec("log", "--pretty=format:%s - %Cred%h%Creset - %C(bold blue)%an%Creset - %Cgreen%cr%Creset", "-n", "1", branch)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return output
}

func (c Commands) IsBehindRemote(branch string) bool {
	output, err := c.exec("status", "-sb", branch)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return strings.Contains(output, "behind")
}

func (c Commands) Fetch() {
	_, err := c.exec("fetch")
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commands) pushBranch(branchName string) {
	fmt.Println("Pushing", color.Yellow(branchName), "...")
	_, err := c.exec("push")
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commands) gitClean() bool {
	output, err := c.exec("status", "--porcelain")
	if err != nil {
		return false
	}
	return len(output) == 0
}

func (c Commands) merge(currentBranch string, toMerge string) error {
	output, err := c.exec(
		"merge",
		toMerge,
		"--no-squash",
		"--commit",
		"-m",
		"Merge branch "+toMerge+" into "+currentBranch+" (gostacking)",
	)
	if err != nil {
		fmt.Println("Error merging:", output)
		return err
	}
	return nil
}
