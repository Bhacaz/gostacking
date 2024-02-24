/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Bhacaz/gostacking/internal/stack"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [branch]",
	Short: "Add a branch to the current stack. If no branch is given, add the current branch.",
	Long: `Add a branch to the current stack.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
		    stack.Add("")
		} else {
            stack.Add(args[0])
        }
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
