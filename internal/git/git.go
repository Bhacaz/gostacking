package git

import (
	"os/exec"
	"strings"
)

type InterfaceGitExecutor interface {
	ExecCommand(command string) (string, error)
}

type Executor struct{}

func (gee Executor) ExecCommand(command string) (string, error) {
	cmdArgs := strings.Fields(command)
	cmd := exec.Command("git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	result := strings.TrimSuffix(string(output), "\n")

	if err != nil {
		//fmt.Println("Git command err:", result)
		return "", err
	}

	return result, nil
}

type Commands struct {
	executor InterfaceGitExecutor
}

func Cmd() Commands {
	return Commands{
		executor: Executor{},
	}
}

func (gc Commands) CurrentBranchName() (string, error) {
	currentBranch, err := gc.executor.ExecCommand("rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", err
	}
	return currentBranch, nil
}

func (gc Commands) BranchExists(branchName string) bool {
	_, err := gc.executor.ExecCommand("rev-parse --verify --quiet \"refs/heads/" + branchName + "\"")
	return err == nil
}

//
// func SyncBranches(branches []string, checkoutBranchEnd string, push bool) {
//     // Return if contains unstaged changes
//     if !gitClean() {
//         fmt.Println("Unstaged changes. Please commit or stash them.")
//         return
//     }
//
//     fmt.Println("Fetching...")
//     _, err := executeGitCommand("fetch")
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//
//     for i, branch := range branches {
//         fmt.Println("Checkout to", branch)
//         _, err := executeGitCommand("checkout " + branch)
//         if err != nil {
//             fmt.Println(err)
//             break
//         }
//
//         fmt.Println("Pulling", branch, "...")
//         _, err = executeGitCommand("pull")
//         if err != nil {
// //             fmt.Println(err)
//         }
//
//         // Nothing to merge on first branch
//         if i == 0 {
//             if push {
//                 pushBranch(branch)
//             }
//             continue
//         }
//         toMerge := branches[i - 1]
//         fmt.Println("Merging", toMerge, "->", branch)
//         err = executeGitMerge(branch, toMerge)
//         if err != nil {
//             fmt.Println(err)
//             break
//         }
//         if push {
//             pushBranch(branch)
//         }
//     }
//     _, err = executeGitCommand("checkout " + checkoutBranchEnd)
//     if err != nil {
//         fmt.Println(err)
//     }
// }
//
// func pushBranch(branchName string) {
//    fmt.Println("Pushing", branchName, "...")
//    _, err := executeGitCommand("push")
//    if err != nil {
//        fmt.Println(err)
//    }
// }
//
// func Checkout(branchName string) {
//     _, err := executeGitCommand("checkout " + branchName)
//     if err != nil {
//         fmt.Println(err)
//     }
// }
//
// func gitClean() bool {
//     output, err := executeGitCommand("status --porcelain")
//     if err != nil {
//         fmt.Println(err)
//         return false
//     }
//     return len(output) == 0
// }
//
// func executeGitMerge(currentBranch string, toMerge string) error {
//     cmd := exec.Command("git", "merge", toMerge, "--no-squash", "--commit", "-m", "Merge branch " + toMerge + " into " + currentBranch + " (gostacking)")
//     output, err := cmd.CombinedOutput()
//     if err != nil {
//         fmt.Println("Error merging:", string(output))
//         return err
//     }
//     return nil
// }
