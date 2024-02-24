/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "strconv"
	"github.com/spf13/cobra"
	"github.com/Bhacaz/gostacking/stack"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [stack name or number]",
	Short: "Change the current stack.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
        if n, err := strconv.Atoi(args[0]); err == nil {
            stack.SwitchByNumber(n)
        } else {
            stack.SwitchByName(args[0])
        }
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
