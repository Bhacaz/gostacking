/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Bhacaz/gostacking/internal/stack"
	"github.com/spf13/cobra"
	"os"
)

var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gostacking",
	Short: "Git stacking with merge.",
	Long: `Git stacking with merge.
Written in go.`,
	Version: "v0.10.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gostacking.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display all Git commands run under the hood")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// 	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func stacksManager() stack.StacksManager {
	return stack.NewManager(Verbose)
}
