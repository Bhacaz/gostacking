package git

import (
    "fmt"
    "gopkg.in/src-d/go-git.v4"
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
