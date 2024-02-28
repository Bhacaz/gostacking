/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Bhacaz/gostacking/internal/stack"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Merge all branches into the others",
	Long: `Merge all branches into the others.
This command will merge all branches into the others, starting from the bottom of the stack.
The current git status must be clean before running this command.
Each branch will be pulled to prevent conflict with the remote.`,
	Run: func(cmd *cobra.Command, args []string) {
		pushValue, _ := cmd.Flags().GetBool("push")
		stack.Manager().Sync(pushValue)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	syncCmd.Flags().BoolP("push", "p", false, "Push commits after syncing.")
}
