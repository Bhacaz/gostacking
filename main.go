/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
    "fmt"
    "github.com/Bhacaz/gostacking/internal/git"
    )

func main() {
    gitCmd := git.GitCmd()
    fmt.Println(gitCmd.CurrentBranchName())
}
