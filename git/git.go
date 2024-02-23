package git

import (
    "fmt"
    "os/exec"
    "strings"
    "gopkg.in/src-d/go-git.v4"
)

func CurrentBranchName() string {
    r, err := git.PlainOpen(".")
    if err != nil {
        fmt.Println(err)
    }

    ref, err := r.Head()
    if err != nil {
        fmt.Println(err)
    }

    return ref.Name().Short()
}

func BranchExists(branchName string) bool {
    r, err := git.PlainOpen(".")
    if err != nil {
        fmt.Println(err)
    }

    refs, err := r.Branches()
    if err != nil {
        fmt.Println(err)
    }

    for {
        ref, err := refs.Next()
        if err != nil {
            break
        }
        if ref.Name().Short() == branchName {
            return true
        }
    }
    return false
}

func SyncBranches(branches []string, checkoutBranchEnd string, push bool) {
    // Return if contains unstaged changes
    if !gitClean() {
        fmt.Println("Unstaged changes. Please commit or stash them.")
        return
    }

    fmt.Println("Fetching...\n")
    _, err := executeGitCommand("fetch")
    if err != nil {
        fmt.Println(err)
        return
    }

    for i, branch := range branches {
        fmt.Println("Checkout to", branch)
        _, err := executeGitCommand("checkout " + branch)
        if err != nil {
            fmt.Println(err)
            break
        }

        fmt.Println("Pulling", branch, "...")
        _, err = executeGitCommand("pull")
        if err != nil {
//             fmt.Println(err)
        }

        // Nothing to merge on first branch
        if i == 0 {
            if push {
                pushBranch(branch)
            }
            continue
        }
        toMerge := branches[i - 1]
        fmt.Println("Merging", toMerge, "->", branch)
        err = executeGitMerge(branch, toMerge)
        if err != nil {
            fmt.Println(err)
            break
        }
        if push {
            pushBranch(branch)
        }
    }
    _, err = executeGitCommand("checkout " + checkoutBranchEnd)
    if err != nil {
        fmt.Println(err)
    }
}

func pushBranch(branchName string) {
   fmt.Println("Pushing", branch, "...")
   _, err = executeGitCommand("push")
   if err != nil {
       fmt.Println(err)
   }
}

func Checkout(branchName string) {
    _, err := executeGitCommand("checkout " + branchName)
    if err != nil {
        fmt.Println(err)
    }
}

func gitClean() bool {
    output, err := executeGitCommand("status --porcelain")
    if err != nil {
        fmt.Println(err)
        return false
    }
    return len(output) == 0
}

func executeGitMerge(currentBranch string, toMerge string) error {
    cmd := exec.Command("git", "merge", toMerge, "--no-squash", "--commit", "-m", "Merge branch " + toMerge + " into " + currentBranch + " (gostacking)")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error merging:", string(output))
        return err
    }
    return nil
}

func executeGitCommand(command string) (string, error) {
    cmdArgs := strings.Fields(command)
    cmd := exec.Command("git", cmdArgs...)
    output, err := cmd.CombinedOutput()
    if err != nil {
//         fmt.Println("Command err:", string(output))
        return "", err
    }
    return string(output), nil
}
