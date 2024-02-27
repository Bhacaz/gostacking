/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Bhacaz/gostacking/internal/stack"
	"github.com/spf13/cobra"
	"strconv"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a branch from the current stack. (Branch name or number)",
	Run: func(cmd *cobra.Command, args []string) {
		if n, err := strconv.Atoi(args[0]); err == nil {
			stack.Manager().RemoveByNumber(n)
		} else {
			stack.Manager().RemoveByName(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
