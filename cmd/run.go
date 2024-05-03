/*
Copyright © 2024 OKAREO oss@okareo.com
*/
package cmd

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"

	"bytes"
	"io"
	"net/http"
)

type FlowConfig struct {
	Name        string   `yaml:"name"`
	Project_id  string   `yaml:"project-id"`
	Model_id    string   `yaml:"model-id"`
	Scenario_id string   `yaml:"scenario-id"`
	Type        string   `yaml:"type"`
	Checks      []string `yaml:"checks"`
}

type ModelValues map[string]string // map of model keys to model ids

type Config struct {
	Name      string
	APIKey    string `yaml:"api-key"`
	ProjectID string `yaml:"project-id"`
	Language  string `yaml:"language"`
	ModelKeys struct {
		Values ModelValues `yaml:"values"`
	}
	Run struct {
		Flows struct {
			FilePattern string        `yaml:"file-pattern"`
			FlowConfigs []*FlowConfig `yaml:"configs"`
		}
	}
}

type TestRun struct {
	ID                 string                   `json:"id"`
	ProjectID          string                   `json:"project_id"`
	MutID              string                   `json:"mut_id"`
	ScenarioSetID      string                   `json:"scenario_set_id"`
	Name               string                   `json:"name"`
	Tags               []string                 `json:"tags"`
	Type               string                   `json:"type"`
	StartTime          string                   `json:"start_time"`
	EndTime            string                   `json:"end_time"`
	TestDataPointCount int                      `json:"test_data_point_count"`
	ModelMetrics       map[string](interface{}) `json:"model_metrics"`
	ErrorMatrix        map[string](interface{}) `json:"error_matrix"`
	AppLink            string                   `json:"app_link"`
}

type Model struct {
	ID             string                 `json:"id"`
	ProjectID      string                 `json:"project_id"`
	Name           string                 `json:"name"`
	Models         map[string]interface{} `json:"models"`
	Tags           []string               `json:"tags"`
	TimeCreated    string                 `json:"time_created"`
	DatapointCount int                    `json:"datapoint_count"`
	AppLink        string                 `json:"app_link"`
	Warning        string                 `json:"warning"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get_endpoint() string {
	endpoint := os.Getenv("OKAREO_BASE_URL")
	if endpoint == "" {
		return "https://api.okareo.com"
	}
	return endpoint
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Okareo CLI command to run workflows.",
	Long:  `Okareo CLI 'runs' can include multiple flows that perform a variety of tasks from scenario generation to model evaluation.`,
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
		var filePattern string = config.Run.Flows.FilePattern
		var language string = config.Language
		var okareoAPIKey string = tradeForEnvValue(config.APIKey)
		var projectId string = tradeForEnvValue(config.ProjectID)
		var runScripts bool = false
		var runConfigFlows bool = (len(config.Run.Flows.FlowConfigs) > 0)
		if language != "" {
			if strings.ToLower(language) == "python" || strings.ToLower(language) == "py" {
				language = "python"
				runScripts = true
			} else if strings.ToLower(language) == "typescript" || strings.ToLower(language) == "ts" {
				language = "typescript"
				runScripts = true
			}
		}

		if !runConfigFlows && !runScripts {
			fmt.Println("No flows or scripts to run.")
			return
		}

		if runConfigFlows {
			for i := 0; i < len(config.Run.Flows.FlowConfigs); i++ {
				model := get_model(okareoAPIKey, config.Run.Flows.FlowConfigs[i].Model_id, isDebug)
				model_type := ""
				for model_type = range model.Models {
					break // just pulling the key to see what kind of model this is
				}
				model_key := ""
				if model_type == "openai" {
					model_key = os.Getenv("OPENAI_API_KEY")
				}
				project_id := model.ProjectID
				config.Run.Flows.FlowConfigs[i].Project_id = project_id
				testrun := run_config_test(okareoAPIKey, model_type, model_key, config.Run.Flows.FlowConfigs[i], isDebug)

				fmt.Println("Completed: " + testrun.Name)
				fmt.Println("ID: " + testrun.ID)
				fmt.Println("Link: " + testrun.AppLink)
				fmt.Println("-----")
			}
		}

		if runScripts {
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
				var dist_folder string = "./.okareo/dist/"

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

				config_entries, err := os.ReadDir(flows_folder + "/config/")
				if err != nil {
					log.Fatal(err)
				}
				for _, e := range config_entries {
					match, _ := regexp.MatchString(filePattern+"$", e.Name())
					var distFile string = strings.Split(e.Name(), ".")[0] + ".js"
					if isDebug {
						fmt.Println("Match file:", e.Name(), match)
					}
					if match {
						fmt.Println("Running .okareo/flows/config/" + e.Name())
						doJSScript(dist_folder+"config/"+distFile, okareoAPIKey, projectId, run_name, isDebug)
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
		}
	},
}

func get_model(api_token string, model_id string, isDebug bool) *Model {
	endpoint := get_endpoint()
	url := endpoint + "/v0/models_under_test/" + model_id
	var client http.Client
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("api-key", api_token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("error", err)
		fmt.Print(err)
	}

	model := &Model{}
	derr := json.NewDecoder(resp.Body).Decode(model)
	if derr != nil {
		println("Error decoding model.", derr)
		panic(derr)
	}
	if isDebug {
		println("Model ID: ", model.ID)
	}

	if resp.StatusCode != http.StatusCreated {
		panic(resp.Status)
	}
	return model
}

func run_config_test(api_token string, model_type string, model_key string, flow *FlowConfig, isDebug bool) *TestRun {
	endpoint := get_endpoint()
	url := endpoint + "/v0/test_run"
	var client http.Client
	checks := ""
	if len(flow.Checks) > 0 {
		checks = `,
		"checks": [`
		for i := 0; i < len(flow.Checks); i++ {
			if i > 0 {
				checks += `,`
			}
			checks += `"` + flow.Checks[i] + `"`
		}
		checks += `]`
	}

	body_str := `{
		"name":"` + flow.Name + `",
		"project_id":"` + flow.Project_id + `",
		"scenario_id":"` + flow.Scenario_id + `",
		"mut_id":"` + flow.Model_id + `",
		"type":"` + flow.Type + `",
		"calculate_metrics": "true",
		"api_keys": {
			"` + model_type + `": "` + model_key + `"
		}` + checks + `
	}`
	body := []byte(body_str)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

	req.Header.Add("api-key", api_token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request. ", err)
		fmt.Print(err)
	}

	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
		log.Fatal(err)
	}
	resultBodyString := string(resultBody)
	if isDebug {
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, resultBody, "", "\t")
		if error != nil {
			log.Println("JSON parse error: ", error)
			return nil
		}
		log.Println("Test Run:", prettyJSON.String())
	}
	testrun := jsonTestDecoder(resultBodyString)
	return testrun
}

func jsonTestDecoder(body string) *TestRun {

	testrun := &TestRun{}
	dec := json.NewDecoder(strings.NewReader(body))

	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if dec.More() && t == "id" {
			v, v_err := dec.Token()
			if v_err == io.EOF {
				break
			}
			testrun.ID = v.(string)
		}
		if dec.More() && t == "app_link" {
			v, v_err := dec.Token()
			if v_err == io.EOF {
				break
			}
			testrun.AppLink = v.(string)
		}
		if dec.More() && t == "name" {
			v, v_err := dec.Token()
			if v_err == io.EOF {
				break
			}
			testrun.Name = v.(string)
		}
	}

	return testrun
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
			"okareo-ts-sdk": "latest",
			"typescript": "^5.4.2"
		},
		"scripts": {
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
	cmd.Stderr = os.Stderr

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
			"okareo-ts-sdk": "latest"
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
	cmd.Stderr = os.Stderr

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
	runCmd.PersistentFlags().BoolP("debug", "d", false, "See additional stdout to debug your flows.")
	runCmd.PersistentFlags().BoolP("upgrade", "u", false, "Force upgrade to the latest Okareo library. Currently only supported for python.")
	runCmd.PersistentFlags().BoolP("latest", "l", true, "Install the latest version of Okareo. False will require you to maintain okareo yourself.")
}
