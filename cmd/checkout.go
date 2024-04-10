/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"strconv"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout [branch or number]",
	Short: "Checkout a branch from a stack",
	Long: `Checkout a branch from a stack.
If a number is given, checkout the branch by its number in the stack (see status command).
If a name is given, checkout the branch by its name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if n, errParse := strconv.Atoi(args[0]); errParse == nil {
			err = stacksManager().CheckoutByNumber(n)
		} else {
			err = stacksManager().CheckoutByName(args[0])
		}
		return err
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return stacksManager().ListBranchesForCompletion(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
