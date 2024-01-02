package git

import (
    "fmt"
    "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/plumbing"
)

func CurrentBranchName() string {
    fmt.Println("Git module currentBranchName called")
    // read the current branch name from the git repo
    // return the current branch name

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
    fmt.Println("Git module syncBranches called")
    for i, branch := range branches {
        if i == 0 {
            continue
        }
        toMerge := branches[i - 1]
        fmt.Println("Merging", toMerge, "into", branch)
        err := w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)})
        err = w.Pull(&git.PullOptions{
                    RemoteName: ".",
                    ReferenceName: plumbing.NewBranchReferenceName(toMerge),
                })
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
