/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "okareo-cli",
	Short: "Use Okaero to evaluate your use of AI/ML in your application.",
	Long: `The Okareo CLI is a tool to help you evaluate your use of AI/ML in your application:

Create .py file in .okareo/validations to contiuosly unit test your AI/ML model.`,
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.okareo-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
