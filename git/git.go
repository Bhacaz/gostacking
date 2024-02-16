package git

import (
    "fmt"
    "os/exec"
    "strings"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
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

func SyncBranches(branches []string, checkoutBranchEnd string) {
    r, _ := git.PlainOpen(".")
    w, _ := r.Worktree()
    // Return if contains unstaged changes
    if !gitClean() {
        fmt.Println("Unstaged changes. Please commit or stash them.")
        return
    }

    fmt.Println("Fetching...")
    _, err := executeGitCommand("fetch")
    if err != nil {
        fmt.Println(err)
        return
    }

    for i, branch := range branches {
        fmt.Println("Checkout to ", branch)
        _, err := executeGitCommand("checkout " + branch)
        if err != nil {
            fmt.Println(err)
            break
        }

        fmt.Println("Pulling ", branch)
        _, err = executeGitCommand("pull")
        if err != nil {
            fmt.Println(err)
            break
        }

        // Nothing to merge on first branch
        if i == 0 {
            continue
        }
        toMerge := branches[i - 1]
        fmt.Println("Merging", toMerge, "into", branch)
        err = executeGitMerge(branch, toMerge)
        if err != nil {
            fmt.Println(err)
            break
        }
    }
    err = w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(checkoutBranchEnd)})
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
    cmd := exec.Command("git", "merge", toMerge, "--no-squash", "--commit", "-m", "\"gostacking - Merge " + toMerge + " into " + currentBranch + "\"")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("Command output:", string(output))
        return err
    }
    return nil
}

func executeGitCommand(command string) (string, error) {
    cmdArgs := strings.Fields(command)
    cmd := exec.Command("git", cmdArgs...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("Command err:", string(output))
        return "", err
    }
//     fmt.Println("Command output:", string(output))
    return string(output), nil
}
