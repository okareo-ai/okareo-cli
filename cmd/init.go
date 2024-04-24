/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a default .okareo structure in the current directory",
	Long:  `Creates a default .okareo structure in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		isDebug, _ := cmd.Flags().GetBool("debug")
		isForce, _ := cmd.Flags().GetBool("force")
		language, _ := cmd.Flags().GetString("language")

		var okareo_folder = "./.okareo"
		var init_file = "config.yml"
		var init_file_path = okareo_folder + "/" + init_file
		var flow_folder = "flows"
		var flow_folder_path = okareo_folder + "/" + flow_folder
		var flow_example_name = "example"
		var flow_example_path = flow_folder_path + "/" + flow_example_name

		if dirExists(okareo_folder) && fileExists(init_file_path) {
			if isForce {
				if fileExists(init_file_path) { //redundant check
					p_err := os.Remove(init_file_path)
					if p_err != nil {
						fmt.Println("ERROR", init_file_path)
						fmt.Println(p_err)
					}
					if isDebug {
						fmt.Println("Removed", init_file_path)
					}
				}
			} else {
				fmt.Println("Your Okareo instance is already initialized. Use --force to overwrite. INFO: This will not delete existing flows.")
				return
			}
		}
		config := []byte(``)
		example_flow := []byte(``)

		if language == "" {
			config = []byte(`name: CLI Evaluation 
api-key: ${OKAREO_API_KEY}
run:
  flows:
    configs:
#      - name: "Example Flow"
#        model-id: "MODEL_ID"
#        scenario-id: "SCENARIO_ID"
#        type: "NL_GENERATION"
#        openai-key: "${OPENAI_API_KEY}"
#        checks:
#          - uniqueness
#          - fluency
`)
		} else if strings.ToLower(language) == "python" || strings.ToLower(language) == "py" {
			config = []byte(`name: CLI Evaluation 
api-key: ${OKAREO_API_KEY}
language: "python"
run:
  flows:
    file-pattern: '.*\.py'
`)
			example_flow = getExamplePythonFlow()
			flow_example_path += ".py"

		} else if strings.ToLower(language) == "javascript" || strings.ToLower(language) == "js" {
			config = []byte(`name: CLI Evaluation 
api-key: ${OKAREO_API_KEY}
language: "javascript"
run:
  flows:
    file-pattern: '.*\.js'
`)
			example_flow = getExampleJavascriptFlow()
			flow_example_path += ".js"

		} else if strings.ToLower(language) == "typescript" || strings.ToLower(language) == "ts" {
			config = []byte(`name: CLI Evaluation 
api-key: ${OKAREO_API_KEY}
language: "typescript"
run:
  flows:
    file-pattern: '.*\.ts'
`)
			example_flow = getExampleTypescriptFlow()
			flow_example_path += ".ts"
		}

		_, err_okareo_folder := os.Stat(okareo_folder)
		if os.IsNotExist(err_okareo_folder) {
			f_err := os.Mkdir(okareo_folder, 0777)
			check(f_err)
		}

		_, err_flow_folder := os.Stat(flow_folder_path)
		if os.IsNotExist(err_flow_folder) {
			ff_err := os.Mkdir(flow_folder_path, 0777)
			check(ff_err)
		}

		_, err_pkg := os.Stat(init_file_path)
		if os.IsNotExist(err_pkg) {
			f_err := os.WriteFile(init_file_path, config, 0777)
			check(f_err)
		}

		_, err_flow := os.Stat(flow_folder_path)
		if os.IsNotExist(err_flow) {
			fff_err := os.Mkdir(flow_folder_path, 0777)
			check(fff_err)
		}

		if language != "" {
			_, err_example := os.Stat(flow_example_path)
			if os.IsNotExist(err_example) {
				f_err := os.WriteFile(flow_example_path, example_flow, 0777)
				check(f_err)
			}
		}

	},
}

func getExamplePythonFlow() []byte {
	return []byte(`#!/usr/bin/env python3
import os
from okareo import Okareo
OKAREO_API_KEY = os.environ["OKAREO_API_KEY"]
okareo = Okareo(OKAREO_API_KEY)
print("Python example to get projects")
projects = okareo.get_projects()
print('Project List: ', projects)
`)
}

func getExampleTypescriptFlow() []byte {
	return []byte(`import { Okareo } from 'okareo-ts-sdk';
const okareo = new Okareo({api_key:process.env.OKAREO_API_KEY });
const main = async () => {
	try {
		console.log("Typescript example to get projects");
		const projects = await okareo.getProjects();
		console.log("List Projects:", projects);
	} catch (error) {
		console.error(error);
	}
}
main();`)
}

func getExampleJavascriptFlow() []byte {
	return []byte(`const okareo_sdk = require('okareo-ts-sdk');
const okareo = new okareo_sdk.Okareo({api_key:process.env.OKAREO_API_KEY });
const main = async () => {
	try {
		console.log("Javascript example to get projects");
		const projects = await okareo.getProjects();
		console.log("List Projects:", projects);
	} catch (error) {
		console.error(error);
	}
}
main();`)
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringP("language", "l", "", "The language you want to configure: Python, Javascript, or Typescript.")
	initCmd.PersistentFlags().BoolP("force", "f", false, "Forces the exiting configuration to be overwritten.")
	initCmd.PersistentFlags().BoolP("debug", "d", false, "See additional stdout to debug the init process.")
}
