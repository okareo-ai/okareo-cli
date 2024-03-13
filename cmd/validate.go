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
	Use:   "validate",
	Short: "Okareo CLI command to run validations",
	Long:  `Okareo CLI command to run validations`,
	Run: func(cmd *cobra.Command, args []string) {
		isDebug, _ := cmd.Flags().GetBool("debug")
		useLatest, _ := cmd.Flags().GetBool("latest")
		if useLatest {
			installOkareo(isDebug)
		}
		validationFileFlag, _ := cmd.Flags().GetString("file")
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

		var okareoAPIKey string = tradeForEnvValue(config.APIKey)
		var projectId string = tradeForEnvValue(config.ProjectID)

		var filePattern string = config.Run.Scripts.FilePattern

		if validationFileFlag != "ALL" {
			fmt.Println("./.okareo/validations/" + validationFileFlag)
			doPythonScript("./.okareo/validations/"+validationFileFlag, okareoAPIKey, projectId, run_name)
		} else {
			entries, err := os.ReadDir("./.okareo/validations/")
			if err != nil {
				log.Fatal(err)
			}

			for _, e := range entries {
				match, _ := regexp.MatchString(filePattern+"$", e.Name())
				if isDebug {
					fmt.Println("Match file:", e.Name(), match)
				}
				if match {
					fmt.Println("Running .okareo/validations/" + e.Name())
					doPythonScript("./.okareo/validations/"+e.Name(), okareoAPIKey, projectId, run_name)
				}
			}
		}

	},
}

func installOkareo(debug bool) {
	// create the install file and overwrite if it already exists
	inst_script := []byte("python3 -m pip install \"okareo\"\n")
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

func doPythonScript(filename string, okareoAPIKey string, projectId string, run_name string) {
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

func tradeForEnvValue(envVar string) string {
	if strings.HasPrefix(envVar, "${") && strings.HasSuffix(envVar, "}") {
		return os.Getenv(envVar[2 : len(envVar)-1])
	}
	return envVar
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().String("file", "ALL", "The Okareo validation script you want to run.")
	runCmd.PersistentFlags().String("config", "./.okareo/config.yml", "The Okareo configuration file for the evaluation run.")
	runCmd.PersistentFlags().Bool("debug", false, "See additional stdout to debug your scripts.")
	runCmd.PersistentFlags().Bool("latest", true, "Install the latest version of Okareo. False will require you to maintain okareo yourself.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
