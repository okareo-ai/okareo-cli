package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
)
func isPythonInstalled() bool {
    cmd := exec.Command("python", "--version")
    return cmd.Run() == nil
}

func isPipInstalled() bool {
    cmd := exec.Command("pip", "--version")
    return cmd.Run() == nil
}

func promptPythonInstallation() {
    fmt.Println("Python is not installed. Please install Python 3.7 or later.")
    fmt.Println("Visit https://www.python.org/downloads/ for installation instructions.")
    fmt.Println("After installing Python, please run this command again.")
    os.Exit(1)
}
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start a proxy server using litellm",
	Long:  `Starts a proxy server that can handle LLM requests using litellm's proxy functionality`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		config, _ := cmd.Flags().GetString("config")
		debug, _ := cmd.Flags().GetBool("debug")

		if debug {
			fmt.Println("Debug mode enabled")
			fmt.Println("Checking for Python and pip...")
		}
	
		if !isPythonInstalled() {
			fmt.Println("Python is not installed.")
			promptPythonInstallation()
			return
		} else {
			if debug {
				fmt.Println("Python is installed")
			}
		}
	
		if !isPipInstalled() {
			if debug {
				fmt.Println("pip not found. Attempting to install...")
			}
			installCmd := exec.Command("python", "-m", "ensurepip", "--upgrade")
			if err := installCmd.Run(); err != nil {
				fmt.Printf("Error installing pip: %v\n", err)
				fmt.Println("Please install pip manually and try again.")
				return
			}
			if debug {
				fmt.Println("pip installed successfully")
			}
		} else if debug {
			fmt.Println("pip is already installed")
		}

		if debug {
			fmt.Println("Installing required packages...")
		}
		// Install litellm and required opentelemetry packages
		installCmd := exec.Command("pip", "install",
			"litellm[proxy]==1.53.7",
			"opentelemetry-api==1.27.0",
			"opentelemetry-exporter-otlp==1.27.0",
			"opentelemetry-exporter-otlp-proto-common==1.27.0", 
			"opentelemetry-exporter-otlp-proto-grpc==1.27.0",
			"opentelemetry-exporter-otlp-proto-http==1.27.0",
			"opentelemetry-instrumentation==0.48b0",
			"opentelemetry-instrumentation-asgi==0.48b0", 
			"opentelemetry-instrumentation-sqlalchemy==0.48b0",
			"opentelemetry-proto==1.27.0",
			"opentelemetry-sdk==1.27.0",
			"opentelemetry-semantic-conventions==0.48b0",
			"opentelemetry-util-http==0.48b0",
		)
		if debug {
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr
		} else {
			installCmd.Stdout = nil
			installCmd.Stderr = nil
		}
		if err := installCmd.Run(); err != nil {
			fmt.Printf("Error installing litellm: %v\n", err)
			return
		}

		// Build the litellm command
		cmdArgs := []string{}
		if port != "" {
			cmdArgs = append(cmdArgs, "--port", port)
		} else {
			port = "4000"
			cmdArgs = append(cmdArgs, "--port", port)
		}

		// Get config file path relative to current directory
		userConfigPath := config
		//"./cmd/proxy_config.yaml"
		// Create a temporary config file with default settings
		defaultConfig := []byte(`model_list:
  - model_name: "*" 
    litellm_params:
      model: "*"
litellm_settings:
  callbacks: ["otel"]`)

		tmpConfig, err := os.CreateTemp("", "proxy_config_*.yaml")
		if err != nil {
			return
		}
		defer func() {
			os.Remove(tmpConfig.Name())
		}()

		if _, err := tmpConfig.Write(defaultConfig); err != nil {
			return
		}

		if err := tmpConfig.Close(); err != nil {
			return
		}

		defaultConfigPath := tmpConfig.Name()

		if config != "" {
			cmdArgs = append([]string{"--config", userConfigPath}, cmdArgs...)
		} else {
			cmdArgs = append([]string{"--config", defaultConfigPath}, cmdArgs...)
		}
		litellmCmd := exec.Command("litellm", cmdArgs...)
		// Get existing env and add OTEL vars
		env := os.Environ()
		filteredEnv := make([]string, 0, len(env))
        for _, e := range env {
            if !strings.HasPrefix(e, "DATABASE_URL=") {
                filteredEnv = append(filteredEnv, e)
            }
        }
        env = filteredEnv
		okareoApiKey := os.Getenv("OKAREO_API_KEY")
		dev, _ := cmd.Flags().GetBool("dev")

		if okareoApiKey != "" {
			if dev {
				env = append(env, "OTEL_ENDPOINT=http://localhost:8000/v0/traces")
			} else {
				env = append(env, "OTEL_ENDPOINT=https://api.okareo.com/v0/traces")
			}
			env = append(env, "OTEL_HEADERS=api-key="+okareoApiKey)
			env = append(env, "OTEL_EXPORTER=otlp_http")
		}
		litellmCmd.Env = env
		if debug {
			litellmCmd.Stdout = os.Stdout
			litellmCmd.Stderr = os.Stderr
		} else {
			litellmCmd.Stdout = nil
			litellmCmd.Stderr = nil
		}

		fmt.Printf("Starting proxy on port %s\n", port)
		if err := litellmCmd.Run(); err != nil {
			fmt.Printf("Error running proxy: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Add proxy-specific flags
	proxyCmd.Flags().StringP("port", "p", "4000", "Port to run the proxy server on")
	proxyCmd.Flags().StringP("host", "H", "0.0.0.0", "Host to run the proxy server on")
	proxyCmd.Flags().StringP("model", "m", "", "Model to use (e.g., gpt-3.5-turbo, claude-2)")
	proxyCmd.Flags().BoolP("debug", "d", false, "Enable debug mode")
	proxyCmd.Flags().BoolP("dev", "", false, "Use local development endpoint for traces")
	proxyCmd.Flags().StringP("config", "c", "", "Path to config file")
}