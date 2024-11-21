package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Start a proxy server using litellm",
	Long:  `Starts a proxy server that can handle LLM requests using litellm's proxy functionality`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		// Install litellm if not already installed
		installCmd := exec.Command("pip", "install", "litellm[proxy]==1.49.2")
		installCmd.Stdout = nil
		installCmd.Stderr = nil
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
		defaultConfigPath := "./cmd/proxy_config.yaml"
		
		// Only use default config if user hasn't provided their own
		if len(args) == 0 {
			// Check if default config file exists and add to args if it does
			if _, err := os.Stat(defaultConfigPath); !os.IsNotExist(err) {
				cmdArgs = append([]string{"--config", defaultConfigPath}, cmdArgs...)
			}
		} else {
			// Use user provided config
			cmdArgs = append([]string{"--config", args[0]}, cmdArgs...)
		}
		
		litellmCmd := exec.Command("litellm", cmdArgs...)
		
		// Get existing env and add OTEL vars
		env := os.Environ()
		okareoApiKey := os.Getenv("OKAREO_API_KEY")
		
		if okareoApiKey != "" {
			env = append(env, "OTEL_ENDPOINT=http://localhost:8000/v0/traces")
			env = append(env, "OTEL_HEADERS=api-key="+okareoApiKey)
		}

		litellmCmd.Env = env
		litellmCmd.Stdout = nil
		litellmCmd.Stderr = nil

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
}