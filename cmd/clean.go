/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Okareo CLI command to clean .okareo directory",
	Long:  `Okareo CLI command to clean .okareo directory`,
	Run: func(cmd *cobra.Command, args []string) {
		isDebug, _ := cmd.Flags().GetBool("debug")
		force, _ := cmd.Flags().GetBool("force")

		var config_file_path string = "./.okareo/config.yml"
		var pkg_file_path = "./.okareo/package.json"
		var pkg_lock_file_path = "./.okareo/package-lock.json"
		var tsconfig_file_path = "./.okareo/tsconfig.json"
		var install_file_path = "./.okareo/install.sh"
		var node_modules_dir_path = "./.okareo/node_modules"
		var dist_dir_path = "./.okareo/dist"
		var reports_dir_path = "./.okareo/reports"

		config := Config{}

		config_file, read_err := os.ReadFile(config_file_path)
		if (read_err != nil) && (read_err.Error() == "open ./.okareo/config.yaml: no such file or directory") {
			fmt.Println("A config.yml was not found. Continuing with default settings.")
		} else {
			err := yaml.Unmarshal([]byte(config_file), &config)
			if err != nil {
				// make errors topical and friendly
				log.Fatalf("error: %v", err)
				return
			}
		}

		if fileExists(install_file_path) {
			py_err := os.RemoveAll(install_file_path)
			if py_err != nil {
				fmt.Println(py_err)
			}
			if isDebug {
				fmt.Println("Removed", install_file_path)
			}
		}

		//if strings.ToLower(language) == "js" || strings.ToLower(language) == "javascript" || strings.ToLower(language) == "ts" || strings.ToLower(language) == "typescript" || force {

		if fileExists(pkg_file_path) {
			if force {
				p_err := os.Remove(pkg_file_path)
				if p_err != nil {
					fmt.Println("ERROR", pkg_file_path)
					fmt.Println(p_err)
				}
				if isDebug {
					fmt.Println("Removed", pkg_file_path)

				}
			} else {
				fmt.Println("Skipping", pkg_file_path, "use --force ALL to remove")
			}
		}

		if fileExists(pkg_lock_file_path) {
			p_err := os.Remove(pkg_lock_file_path)
			if p_err != nil {
				fmt.Println(p_err)
			}
			if isDebug {
				fmt.Println("Removed", pkg_lock_file_path)
			}
		}
		if fileExists(tsconfig_file_path) {
			tsc_err := os.Remove(tsconfig_file_path)
			if tsc_err != nil {
				fmt.Println(tsc_err)
			}
			if isDebug {
				fmt.Println("Removed", tsconfig_file_path)
			}
		}
		if dirExists(node_modules_dir_path) {
			nm_err := os.RemoveAll(node_modules_dir_path)
			if nm_err != nil {
				fmt.Println(nm_err)
			}
			if isDebug {
				fmt.Println("Removed", node_modules_dir_path)
			}
		}
		if dirExists(dist_dir_path) {
			d_err := os.RemoveAll(dist_dir_path)
			if d_err != nil {
				fmt.Println(d_err)
			}
			if isDebug {
				fmt.Println("Removed", dist_dir_path)
			}
		}
		if dirExists(reports_dir_path) {
			d_err := os.RemoveAll(reports_dir_path)
			if d_err != nil {
				fmt.Println(d_err)
			}
			if isDebug {
				fmt.Println("Removed", reports_dir_path)
			}
		}
		//}

	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.PersistentFlags().BoolP("debug", "d", false, "See additional stdout to debug your scripts.")
	cleanCmd.PersistentFlags().BoolP("force", "f", false, "WARNING: This will remove all configuration files in the .okareo directory. Use with caution.")
}
