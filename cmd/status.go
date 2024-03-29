/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get current stack",
	Long: `Get current stack.
Show the current stack and the current branch.
Branches out of sync with the previous branch are marked with a star (*).
Add the --log flag to show the last commit log for each branch in the stack.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showLogValue, _ := cmd.Flags().GetBool("log")
		return stacksManager().CurrentStackStatus(showLogValue)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	statusCmd.Flags().BoolP("log", "l", false, "Show last commit log for each branch in the stack.")
}
