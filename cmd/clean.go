/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Okareo CLI command to clean .okareo directory",
	Long:  `Okareo CLI command to clean .okareo directory`,
	Run: func(cmd *cobra.Command, args []string) {
		isDebug, _ := cmd.Flags().GetBool("debug")
		files, _ := cmd.Flags().GetString("files")
		configFileFlag, _ := cmd.Flags().GetString("config")

		//var pkg_file = "./.okareo/package.json"
		var pkg_lock_file = "./.okareo/package-lock.json"
		var tsconfig_file = "./.okareo/tsconfig.json"
		var install_file = "./.okareo/install.sh"
		var node_modules_dir = "./.okareo/node_modules"
		var dist_dir = "./.okareo/dist"

		config_file, read_err := os.ReadFile(configFileFlag)
		check(read_err)

		config := Config{}
		err := yaml.Unmarshal([]byte(config_file), &config)
		if err != nil {
			// make errors topical and friendly
			log.Fatalf("error: %v", err)
			return
		}

		var language string = "python"
		if config.Language != "" {
			language = config.Language
		}
		fmt.Println("Cleaning", language, "files")
		if language == "python" || files == "ALL" {
			if fileExists(install_file) {
				py_err := os.RemoveAll(install_file)
				if py_err != nil {
					fmt.Println(py_err)
				}
				if isDebug {
					fmt.Println("Removed", install_file)
				}
			}
		}
		if strings.ToLower(language) == "js" || strings.ToLower(language) == "javascript" || strings.ToLower(language) == "ts" || strings.ToLower(language) == "typescript" || files == "ALL" {
			/*
				if fileExists(pkg_file) {
					p_err := os.Remove(pkg_file)
					if p_err != nil {
						fmt.Println("ERROR", pkg_file)
						fmt.Println(p_err)
					}
					if isDebug {
						fmt.Println("Removed", pkg_file)
					}
				}
			*/
			if fileExists(pkg_lock_file) {
				p_err := os.Remove(pkg_lock_file)
				if p_err != nil {
					fmt.Println(p_err)
				}
				if isDebug {
					fmt.Println("Removed", pkg_lock_file)
				}
			}
			if fileExists(tsconfig_file) {
				if strings.ToLower(language) == "ts" || strings.ToLower(language) == "typescript" || files == "ALL" {
					tsc_err := os.Remove(tsconfig_file)
					if tsc_err != nil {
						fmt.Println(tsc_err)
					}
					if isDebug {
						fmt.Println("Removed", tsconfig_file)
					}
				}
			}
			if dirExists(node_modules_dir) {
				nm_err := os.RemoveAll(node_modules_dir)
				if nm_err != nil {
					fmt.Println(nm_err)
				}
				if isDebug {
					fmt.Println("Removed", node_modules_dir)
				}
			}
			if dirExists(dist_dir) {
				d_err := os.RemoveAll(dist_dir)
				if d_err != nil {
					fmt.Println(d_err)
				}
				if isDebug {
					fmt.Println("Removed", dist_dir)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.PersistentFlags().String("config", "./.okareo/config.yml", "The Okareo configuration file for the evaluation run.")
	cleanCmd.PersistentFlags().Bool("debug", false, "See additional stdout to debug your scripts.")
	cleanCmd.PersistentFlags().String("files", "ALL", "The Okareo flow you want to run.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
