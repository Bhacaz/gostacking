/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Bhacaz/gostacking/internal/stack"
	"github.com/spf13/cobra"
	"strconv"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [stack or number]",
	Short: "Change the current stack",
	Long: `Change the current stack.
If a number is given, switch to the stack by its number in the list of stacks (see list command).
If a name is given, switch to the stack by its name.
If no argument is given, switch to the stack that contains the current branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			stack.Manager().SwitchByName("")
		} else if n, err := strconv.Atoi(args[0]); err == nil {
			stack.Manager().SwitchByNumber(n)
		} else {
			stack.Manager().SwitchByName(args[0])
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return stack.Manager().ListStacksForCompletion(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
