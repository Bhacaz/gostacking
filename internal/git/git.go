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

type InterfaceCommands interface {
	CurrentBranchName() (string, error)
	BranchExists(branchName string) bool
	Checkout(branchName string)
	SyncBranches(branches []string, checkoutBranchEnd string, push bool, mergeHead bool)
	BranchDiff(baseBranch string, branch string) bool
	LastLog(branch string) string
	IsBehindRemote(branch string) bool
	Fetch()
}

//type Commands struct {
//	executor interfaceGitExecutor
//}
//
//func (c Commands) exec(command ...string) (string, error) {
//	return c.executor.execCommand(command...)
//}
//
//func Cmd() Commands {
//	return Commands{
//		executor: executor{},
//	}
//}
//
//func (c Commands) SyncBranches(branches []string, checkoutBranchEnd string, push bool, mergeHead bool) {
//	// Return if contains unstaged changes
//	if !c.gitClean() {
//		fmt.Println("Unstaged changes. Please commit or stash them.")
//		return
//	}
//
//	fmt.Println("Fetching", "...")
//	c.Fetch()
//
//	for i, branch := range branches {
//		fmt.Println("Branch:", color.Yellow(branch))
//		fmt.Println("\tCheckout", "...")
//		_, err := c.exec("checkout", branch)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		fmt.Println("\tPull", "...")
//		_, err = c.exec("pull")
//		if err != nil {
//			fmt.Println("\tNothing to pull")
//		}
//
//		// Nothing to merge on first branch
//		if i == 0 {
//			if mergeHead {
//				err = c.mergeHead(branch)
//				if err != nil {
//					fmt.Println(err)
//					// Stop sync so the user can resolve the conflict
//					return
//				}
//			}
//			if push {
//				c.pushBranch(branch)
//			}
//			continue
//		}
//
//		toMerge := branches[i-1]
//		fmt.Println("\tMerging", color.Yellow(toMerge))
//		err = c.merge(branch, toMerge)
//		if err != nil {
//			fmt.Println(err)
//			// Stop sync so the user can resolve the conflict
//			return
//		}
//		if push {
//			c.pushBranch(branch)
//		}
//	}
//	c.Checkout(checkoutBranchEnd)
//}
//
//func (c Commands) mergeHead(branch string) error {
//	head, err := c.exec("symbolic-ref", "refs/remotes/origin/HEAD", "--short")
//	if err != nil {
//		fmt.Println("Error getting remote HEAD:\n To set it try:", color.Teal("git remote set-head origin main"))
//		return err
//	}
//
//	if !c.BranchDiff(head, branch) {
//		fmt.Println("\tAlready up-to-date with HEAD", color.Yellow(head))
//		return nil
//	}
//
//	fmt.Println("\tMerging HEAD", color.Yellow(head))
//	err = c.merge(branch, head)
//	if err != nil {
//		return err
//	}
//	return nil
//}
