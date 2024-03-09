/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [branch]",
	Short: "Add a branch to the current stack",
	Long: `Add a branch to the current stack.
If no branch is given, add the current branch.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := ""
		if len(args) > 0 {
			branchName = args[0]
		}
		position, _ := cmd.Flags().GetInt("position")
		if position < 0 {
			cmd.PrintErrf("Invalid position: %d. Position must be greater than or equal to 1.", position)
			return nil
		}
		return stacksManager().AddBranch(branchName, position)
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

	// If position is 0 (default value), the branch will be added at the top of the stack.
	addCmd.Flags().IntP("position", "p", 0, "Add the branch at the given position in the stack.")
}
