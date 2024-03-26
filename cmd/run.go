/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name      string
	APIKey    string `yaml:"api-key"`
	ProjectID string `yaml:"project-id"`
	Language  string `yaml:"language"`
	Run       struct {
		Scripts struct {
			FilePattern string `yaml:"file-pattern"`
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Okareo CLI command to run workflows.",
	Long:  `Okareo CLI 'runs' can include multiple scripts that perform a variety of tasks from scenario generation to model evaluation.`,
	Run: func(cmd *cobra.Command, args []string) {
		isDebug, _ := cmd.Flags().GetBool("debug")
		doUpgrade, _ := cmd.Flags().GetBool("upgrade")
		flowFileFlag, _ := cmd.Flags().GetString("file")
		configFileFlag, _ := cmd.Flags().GetString("config")

		config_file, read_err := os.ReadFile(configFileFlag)
		check(read_err)

		config := Config{}
		err := yaml.Unmarshal([]byte(config_file), &config)
		if err != nil {
			// make errors topical and friendly
			log.Fatalf("error: %v", err)
			return
		}

		n := 5
		b := make([]byte, n)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
		s := fmt.Sprintf("%X", b)
		var run_name string = config.Name + "-" + s

		var flows_folder string = "./.okareo/flows/"
		var filePattern string = config.Run.Scripts.FilePattern
		var okareoAPIKey string = tradeForEnvValue(config.APIKey)
		var projectId string = tradeForEnvValue(config.ProjectID)
		var language string = "python"
		if config.Language != "" {
			language = config.Language
		}
		if strings.ToLower(language) == "python" || strings.ToLower(language) == "py" {
			useLatest, _ := cmd.Flags().GetBool("latest")
			if useLatest {
				installOkareoPython(doUpgrade, isDebug)
			}

			entries, err := os.ReadDir(flows_folder)
			if err != nil {
				log.Fatal(err)
			}

			var foundFlow bool = false
			for _, e := range entries {
				if flowFileFlag != "ALL" {
					isFile, _ := regexp.MatchString(flowFileFlag+"(.py)", e.Name())
					if isFile {
						if isDebug {
							fmt.Println("Match file:", e.Name())
						}
						foundFlow = true
						doPythonScript(flows_folder+e.Name(), okareoAPIKey, projectId, run_name, isDebug)
					}
				} else {
					match, _ := regexp.MatchString(filePattern+"$", e.Name())
					if isDebug {
						fmt.Println("Match file:", e.Name(), match)
					}
					if match {
						fmt.Println("Running .okareo/flows/" + e.Name())
						doPythonScript(flows_folder+e.Name(), okareoAPIKey, projectId, run_name, isDebug)
					}
				}
			}
			if flowFileFlag != "ALL" && !foundFlow {
				fmt.Println("Flow not found: " + flowFileFlag)
			}
			//}
		} else if strings.ToLower(language) == "ts" || strings.ToLower(language) == "typescript" {
			installOkareoTypescript(isDebug)
			doTSBuild(isDebug)
			var dist_folder string = "./.okareo/dist/flows/"

			entries, err := os.ReadDir(flows_folder)
			if err != nil {
				log.Fatal(err)
			}

			var foundFlow bool = false
			for _, e := range entries {
				if flowFileFlag != "ALL" {
					isFile, _ := regexp.MatchString(flowFileFlag+"(.ts)", e.Name())
					if isFile {
						if isDebug {
							fmt.Println("Match file:", e.Name())
						}
						foundFlow = true
						var distFile string = strings.Split(e.Name(), ".")[0] + ".js"
						doJSScript(dist_folder+distFile, okareoAPIKey, projectId, run_name, isDebug)
					}
				} else {
					match, _ := regexp.MatchString(filePattern+"$", e.Name())
					var distFile string = strings.Split(e.Name(), ".")[0] + ".js"
					if isDebug {
						fmt.Println("Match file:", e.Name(), match)
					}
					if match {
						fmt.Println("Running .okareo/flows/" + e.Name())
						doJSScript(dist_folder+distFile, okareoAPIKey, projectId, run_name, isDebug)
					}
				}
			}
			if flowFileFlag != "ALL" && !foundFlow {
				fmt.Println("Flow not found: " + flowFileFlag)
			}

		} else if strings.ToLower(language) == "js" || strings.ToLower(language) == "javascript" {
			installOkareoJavascript(isDebug)
			entries, err := os.ReadDir(flows_folder)
			if err != nil {
				log.Fatal(err)
			}

			var foundFlow bool = false
			for _, e := range entries {
				if flowFileFlag != "ALL" {
					isFile, _ := regexp.MatchString(flowFileFlag+"(.js)", e.Name())
					if isFile {
						if isDebug {
							fmt.Println("Match file:", e.Name())
						}
						foundFlow = true
						doJSScript(flows_folder+e.Name(), okareoAPIKey, projectId, run_name, isDebug)
					}
				} else {
					match, _ := regexp.MatchString(filePattern+"$", e.Name())
					if isDebug {
						fmt.Println("Match file:", e.Name(), match)
					}
					if match {
						fmt.Println("Running .okareo/flows/" + e.Name())
						doJSScript(flows_folder+e.Name(), okareoAPIKey, projectId, run_name, isDebug)
					}
				}
			}
			if flowFileFlag != "ALL" && !foundFlow {
				fmt.Println("Flow not found: " + flowFileFlag)
			}
		} else {
			fmt.Println("Language not supported.")
		}
	},
}

func installOkareoPython(doUpgrade bool, debug bool) {
	// create the install file and overwrite if it already exists
	inst_script := []byte("python3 -m pip install \"okareo\"\n")
	if doUpgrade {
		inst_script = []byte("python3 -m pip install --upgrade \"okareo\"\n")
	}
	f_err := os.WriteFile("./.okareo/install.sh", inst_script, 0777)
	check(f_err)

	// run the install script
	cmd := exec.Command("sh", "./.okareo/install.sh")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if debug {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')
		for err == nil {
			fmt.Print(line)
			line, err = reader.ReadString('\n')
		}
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func doPythonScript(filename string, okareoAPIKey string, projectId string, run_name string, isDebug bool) {
	cmd := exec.Command("python3", filename)

	// Setup the environment for the caller
	cmd.Env = os.Environ()
	if run_name != "" {
		cmd.Env = append(cmd.Env, "OKAREO_RUN_ID="+run_name)
	}
	if okareoAPIKey != "" {
		cmd.Env = append(cmd.Env, "OKAREO_API_KEY="+okareoAPIKey)
	}
	if projectId != "" {
		cmd.Env = append(cmd.Env, "PROJECT_ID="+projectId)
	}

	// setup the output handling and call the script
	pipe, err := cmd.StdoutPipe()
	if isDebug {
		cmd.Stderr = os.Stderr
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Println(line)
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func installOkareoTypescript(debug bool) {
	// create the tsconfig file and overwrite if it already exists
	tsconfig_json := []byte(`
	{
		"compilerOptions": {
		  "module": "commonjs",
		  "esModuleInterop": true,
		  "target": "es6",
		  "moduleResolution": "node",
		  "sourceMap": true,
		  "outDir": "dist"
		},
		"lib": ["es2015"]
	}  
	`)
	var tsconfig_file string = "./.okareo/tsconfig.json"
	_, err := os.Stat(tsconfig_file)
	if os.IsNotExist(err) {
		ftsc_err := os.WriteFile(tsconfig_file, tsconfig_json, 0777)
		check(ftsc_err)
	}

	// create the package.json file and overwrite if it already exists
	package_json := []byte(`
	{
		"name": "ts-minimal-ci",
		"version": "0.0.1",
		"description": "Okareo TS Recipe",
		"main": "index.ts",
		"author": "Okareo @ 2024",
		"private": "true",
		"devDependencies": {
			"@types/node": "^20.11.28",
			"okareo-ts-sdk": "^0.0.23",
			"typescript": "^5.4.2"
		},
		"scripts": {
			"start": "tsc && node dist/flows/index.js",
			"build": "tsc"
		}
	}	  
	`)
	var package_file string = "./.okareo/package.json"
	_, err_pkg := os.Stat(package_file)
	if os.IsNotExist(err_pkg) {
		f_err := os.WriteFile(package_file, package_json, 0777)
		check(f_err)
	}

	cmd := exec.Command("npm", "install")
	cmd.Dir = "./.okareo"
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if debug {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')
		for err == nil {
			fmt.Print(line)
			line, err = reader.ReadString('\n')
		}
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}

func doTSBuild(isDebug bool) {
	println("Building typescript flows")
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = "./.okareo"
	// Setup the environment for the caller
	cmd.Env = os.Environ()
	if isDebug {
		cmd.Stderr = os.Stderr
	}

	// setup the output handling and call the script
	pipe, err := cmd.StdoutPipe()
	if isDebug {
		cmd.Stderr = os.Stderr
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func doTSScript(filename string, okareoAPIKey string, projectId string, run_name string, isDebug bool) {
	println("Running", filename)
	cmd := exec.Command("node", filename)
	cmd.Dir = "./.okareo"
	// Setup the environment for the caller
	cmd.Env = os.Environ()
	if isDebug {
		cmd.Stderr = os.Stderr
	}
	if run_name != "" {
		cmd.Env = append(cmd.Env, "OKAREO_RUN_ID="+run_name)
	}
	if okareoAPIKey != "" {
		cmd.Env = append(cmd.Env, "OKAREO_API_KEY="+okareoAPIKey)
	}
	if projectId != "" {
		cmd.Env = append(cmd.Env, "PROJECT_ID="+projectId)
	}

	// setup the output handling and call the script
	pipe, err := cmd.StdoutPipe()
	if isDebug {
		cmd.Stderr = os.Stderr
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func installOkareoJavascript(debug bool) {
	// create the package.json file and overwrite if it already exists
	package_json := []byte(`
	{
		"name": "js-minimal-ci",
		"version": "0.0.1",
		"description": "Okareo JS Recipe",
		"author": "Okareo @ 2024",
		"private": "true",
		"devDependencies": {
			"okareo-ts-sdk": "^0.0.23"
		}
	}	  
	`)
	var package_file string = "./.okareo/package.json"
	_, err_pkg := os.Stat(package_file)
	if os.IsNotExist(err_pkg) {
		f_err := os.WriteFile(package_file, package_json, 0777)
		check(f_err)
	}

	cmd := exec.Command("npm", "install")
	cmd.Dir = "./.okareo"
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if debug {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')
		for err == nil {
			fmt.Print(line)
			line, err = reader.ReadString('\n')
		}
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func doJSScript(filename string, okareoAPIKey string, projectId string, run_name string, isDebug bool) {
	cmd := exec.Command("node", filename)

	// Setup the environment for the caller
	cmd.Env = os.Environ()
	if run_name != "" {
		cmd.Env = append(cmd.Env, "OKAREO_RUN_ID="+run_name)
	}
	if okareoAPIKey != "" {
		cmd.Env = append(cmd.Env, "OKAREO_API_KEY="+okareoAPIKey)
	}
	if projectId != "" {
		cmd.Env = append(cmd.Env, "PROJECT_ID="+projectId)
	}

	// setup the output handling and call the script
	pipe, err := cmd.StdoutPipe()
	if isDebug {
		cmd.Stderr = os.Stderr
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func tradeForEnvValue(envVar string) string {
	if strings.HasPrefix(envVar, "${") && strings.HasSuffix(envVar, "}") {
		return os.Getenv(envVar[2 : len(envVar)-1])
	}
	return envVar
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringP("file", "f", "ALL", "The Okareo flow script you want to run.")
	runCmd.PersistentFlags().StringP("config", "c", "./.okareo/config.yml", "The Okareo configuration file for the evaluation run.")
	runCmd.PersistentFlags().BoolP("debug", "d", false, "See additional stdout to debug your scripts.")
	runCmd.PersistentFlags().BoolP("upgrade", "u", false, "Force upgrade to the latest Okareo library. Currently only supported for python.")
	runCmd.PersistentFlags().BoolP("latest", "l", true, "Install the latest version of Okareo. False will require you to maintain okareo yourself.")
}
