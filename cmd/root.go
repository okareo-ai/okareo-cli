/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "okareo",
	Short: "Use Okaero to evaluate your use of AI/ML in your application.",
	Long: `The Okareo CLI is a tool to help you evaluate your use of AI/ML in your application:
To use the CLI, refer to the docs: https://docs.okareo.com/docs/sdk/cli`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		vFull, _ := cmd.Flags().GetBool("version")
		if vFull {
			fmt.Println("v0.0.7")
		}
	},
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
	rootCmd.PersistentFlags().BoolP("version", "v", true, "The current version of the Okareo CLI")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
