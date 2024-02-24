/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "strconv"
	"github.com/spf13/cobra"
    "github.com/Bhacaz/gostacking/stack"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout [branch name or number]",
	Short: "Checkout a branch from a stack.",
	Run: func(cmd *cobra.Command, args []string) {
		if n, err := strconv.Atoi(args[0]); err == nil {
            stack.CheckoutByNumber(n)
        } else {
            stack.CheckoutByName(args[0])
        }
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
