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
	Use:   "remove [branch or number]",
	Short: "Remove a branch from the current stack",
	Long: `Remove a branch from the current stack.
If a number is given, remove the branch by its number in the stack (see status command).
If a name is given, remove the branch by its name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if n, err := strconv.Atoi(args[0]); err == nil {
			err = stack.NewManager(Verbose).RemoveByNumber(n)
		} else {
			err = stack.NewManager(Verbose).RemoveByName(args[0])
		}
		return err
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
