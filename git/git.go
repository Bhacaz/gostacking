package git

import (
    "fmt"
    "os/exec"
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
    status, _ := w.Status()
    if !status.IsClean() {
        fmt.Println("Unstaged changes. Please commit or stash.")
        return
    }

    for i, branch := range branches {
        fmt.Println("Branch:", branch)
        err := w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)})
        err = w.Pull(&git.PullOptions{})
        fmt.Println("Pulling", branch)
        if err != nil {
            fmt.Println(err, ". Continuing...")
            break
        }

        // Nothing to merge on first branch
        if i == 0 {
            continue
        }
        toMerge := branches[i - 1]
        fmt.Println("Merging", toMerge, "into", branch)
        cmd := exec.Command("git", "pull", ".", toMerge)
        err = cmd.Run()
        if err != nil {
            fmt.Println(err)
            break
        }
    }
    err := w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(checkoutBranchEnd)})
    if err != nil {
        fmt.Println(err)
    }
}
