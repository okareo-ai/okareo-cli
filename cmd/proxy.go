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
		config, _ := cmd.Flags().GetString("config")
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			fmt.Println("Debug mode enabled")
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
			"opentelemetry-instrumentation-fastapi==0.48b0",
			"opentelemetry-instrumentation-sqlalchemy==0.48b0",
			"opentelemetry-proto==1.27.0",
			"opentelemetry-sdk==1.27.0",
			"opentelemetry-semantic-conventions==0.48b0",
			"opentelemetry-util-http==0.48b0",
			"agentops==0.3.14",
			"aiofiles==23.2.1",
			"aiohappyeyeballs==2.4.0",
			"aiohttp==3.10.5",
			"aiosignal==1.3.1",
			"alembic==1.13.3",
			"altair==5.2.0",
			"annotated-types==0.6.0",
			"anthropic==0.20.0",
			"antlr4-python3-runtime==4.8",
			"anyio==4.3.0",
			"appdirs==1.4.4",
			"appnope==0.1.4",
			"APScheduler==3.10.4",
			"arxiv==2.1.3",
			"asgiref==3.7.2",
			"asttokens==2.4.1",
			"async-timeout==5.0.1",
			"attrs==23.2.0",
			"auth0-python==4.7.2",
			"autogen==0.3.1",
			"backoff==2.2.1",
			"bcrypt==4.1.2",
			"beautifulsoup4==4.12.3",
			"bitarray==2.9.2",
			"blis==0.7.11",
			"boto3==1.35.31",
			"botocore==1.35.31",
			"build==1.1.1",
			"CacheControl==0.14.0",
			"cachetools==5.3.3",
			"catalogue==2.0.10",
			"certifi==2024.8.30",
			"cffi==1.16.0",
			"charset-normalizer==3.3.2",
			"chroma-hnswlib==0.7.3",
			"chromadb==0.4.24",
			"cleo==2.1.0",
			"click==8.1.7",
			"cloudpathlib==0.16.0",
			"cloudpickle==2.2.1",
			"cohere==5.11.0",
			"colorama==0.4.6",
			"coloredlogs==15.0.1",
			"comm==0.2.1",
			"confection==0.1.4",
			"contourpy==1.2.0",
			"crashtest==0.4.1",
			"crewai==0.67.1",
			"crewai-tools==0.13.2",
			"cryptography==42.0.8",
			"cycler==0.12.1",
			"cymem==2.0.8",
			"Cython==3.0.9",
			"dataclasses-json==0.6.7",
			"datasets==2.18.0",
			"debugpy==1.8.1",
			"decorator==5.1.1",
			"Deprecated==1.2.14",
			"deprecation==2.1.0",
			"dill==0.3.9",
			"diskcache==5.6.3",
			"distlib==0.3.8",
			"distro==1.9.0",
			"dnspython==2.7.0",
			"docker==7.1.0",
			"docstring_parser==0.16",
			"docx2txt==0.8",
			"dulwich==0.21.7",
			"ed25519==1.5",
			"email_validator==2.2.0",
			"embedchain==0.1.122",
			"en-core-web-lg==3.7.1",
			"en-core-web-sm==3.7.1",
			"evaluate==0.4.1",
			"executing==2.0.1",
			"fairseq==0.12.2",
			"fastapi==0.111.1",
			"fastapi-cli==0.0.5",
			"fastapi-sso==0.10.0",
			"fastavro==1.9.7",
			"fastjsonschema==2.19.1",
			"feedparser==6.0.11",
			"ffmpy==0.3.2",
			"filelock==3.13.1",
			"FLAML==2.2.0",
			"flatbuffers==23.5.26",
			"fonttools==4.49.0",
			"frozendict==2.4.5",
			"frozenlist==1.4.1",
			"fsspec==2024.9.0",
			"gensim==4.3.2",
			"google-api-core==2.18.0",
			"google-auth==2.28.1",
			"google-cloud-aiplatform==1.69.0",
			"google-cloud-bigquery==3.26.0",
			"google-cloud-core==2.4.1",
			"google-cloud-resource-manager==1.12.5",
			"google-cloud-storage==2.16.0",
			"google-crc32c==1.5.0",
			"google-pasta==0.2.0",
			"google-resumable-media==2.7.0",
			"googleapis-common-protos==1.62.0",
			"gptcache==0.1.44",
			"gradio==4.21.0",
			"gradio_client==0.12.0",
			"grpc-google-iam-v1==0.13.1",
			"grpcio==1.66.2",
			"grpcio-status==1.62.3",
			"grpcio-tools==1.62.3",
			"gunicorn==22.0.0",
			"h11==0.14.0",
			"h2==4.1.0",
			"hpack==4.0.0",
			"html5lib==1.1",
			"httpcore==1.0.4",
			"httptools==0.6.1",
			"httpx==0.25.2",
			"httpx-sse==0.4.0",
			"huggingface-hub==0.20.3",
			"humanfriendly==10.0",
			"hydra-core==1.0.7",
			"hyperframe==6.0.1",
			"idna==3.6",
			"importlib-metadata==6.11.0",
			"importlib_resources==6.1.2",
			"inflection==0.5.1",
			"iniconfig==2.0.0",
			"installer==0.7.0",
			"instructor==1.3.3",
			"ipykernel==6.29.2",
			"ipython==8.22.1",
			"ipywidgets==8.1.2",
			"jaraco.classes==3.4.0",
			"jedi==0.19.1",
			"jellyfish==1.0.3",
			"Jinja2==3.1.3",
			"jiter==0.4.2",
			"jmespath==1.0.1",
			"joblib==1.3.2",
			"json_repair==0.25.3",
			"jsonpatch==1.33",
			"jsonpickle==3.3.0",
			"jsonpointer==3.0.0",
			"jsonref==1.1.0",
			"jsonschema==4.23.0",
			"jsonschema-specifications==2023.12.1",
			"jupyter_client==8.6.0",
			"jupyter_core==5.7.1",
			"jupyterlab_widgets==3.0.10",
			"keyring==24.3.1",
			"kiwisolver==1.4.5",
			"kubernetes==29.0.0",
			"lancedb==0.5.7",
			"langchain==0.2.16",
			"langchain-cohere==0.1.9",
			"langchain-community==0.2.17",
			"langchain-core==0.2.41",
			"langchain-experimental==0.0.65",
			"langchain-openai==0.1.25",
			"langchain-text-splitters==0.2.4",
			"langcodes==3.3.0",
			"langsmith==0.1.129",
			"langtrace-python-sdk==3.0.2",
			"Levenshtein==0.25.0",
			"logfire==1.2.0",
			"lxml==5.1.0",
			"Mako==1.3.5",
			"markdown-it-py==3.0.0",
			"MarkupSafe==2.1.5",
			"marshmallow==3.22.0",
			"matplotlib==3.8.3",
			"matplotlib-inline==0.1.6",
			"mdurl==0.1.2",
			"mem0ai==0.1.22",
			"mmh3==4.1.0",
			"mock==4.0.3",
			"monotonic==1.6",
			"more-itertools==10.4.0",
			"mpmath==1.3.0",
			"msgpack==1.0.8",
			"multidict==6.0.5",
			"multiprocess==0.70.17",
			"multitasking==0.0.11",
			"murmurhash==1.0.10",
			"mypy-extensions==1.0.0",
			"nats-py==2.4.0",
			"nbformat==5.9.2",
			"neo4j==5.25.0",
			"nest-asyncio==1.6.0",
			"networkx==3.2.1",
			"nkeys==0.2.0",
			"nltk==3.8.1",
			"nodeenv==1.9.1",
			"numpy==1.26.4",
			"oauthlib==3.2.2",
			"okareo==0.0.75",
			"omegaconf==2.0.6",
			"onnxruntime==1.17.1",
			"openai==1.55.3",
			"openapi-python-client==0.16.0",
			"orjson==3.9.15",
			"outcome==1.3.0.post0",
			"overrides==7.7.0",
			"packaging==23.2",
			"pandas==2.2.0",
			"parameterized==0.9.0",
			"parso==0.8.3",
			"pathos==0.3.3",
			"peewee==3.17.6",
			"pexpect==4.9.0",
			"pillow==10.2.0",
			"pip==24.0",
			"pipdeptree==2.18.1",
			"pkginfo==1.11.1",
			"platformdirs==4.2.0",
			"plotly==5.19.0",
			"pluggy==1.5.0",
			"poetry==1.8.3",
			"poetry-core==1.9.0",
			"poetry-plugin-export==1.8.0",
			"portalocker==2.8.2",
			"posthog==3.6.6",
			"pox==0.3.5",
			"ppft==1.7.6.9",
			"preshed==3.0.9",
			"prompt-toolkit==3.0.43",
			"proto-plus==1.23.0",
			"protobuf==4.25.3",
			"psutil==5.9.8",
			"ptyprocess==0.7.0",
			"pulsar-client==3.4.0",
			"pure-eval==0.2.2",
			"py==1.11.0",
			"pyarrow==15.0.1",
			"pyarrow-hotfix==0.6",
			"pyasn1==0.5.1",
			"pyasn1-modules==0.3.0",
			"pyautogen==0.3.0",
			"pycparser==2.21",
			"pydantic==2.9.2",
			"pydantic_core==2.23.4",
			"pydantic-settings==2.6.0",
			"pydub==0.25.1",
			"Pygments==2.17.2",
			"PyJWT==2.9.0",
			"pylance==0.9.18",
			"PyNaCl==1.5.0",
			"pyparsing==3.1.2",
			"pypdf==4.3.1",
			"PyPDF2==3.0.1",
			"PyPika==0.48.9",
			"pyproject_hooks==1.0.0",
			"pyright==1.1.382.post1",
			"pysbd==0.3.4",
			"PySocks==1.7.1",
			"pytest==8.3.2",
			"pytest-dotenv==0.5.2",
			"python-dateutil==2.8.2",
			"python-dotenv==1.0.1",
			"python-multipart==0.0.9",
			"pytorch-transformers==1.2.0",
			"pytube==15.0.0",
			"pytz==2024.1",
			"pyvis==0.3.2",
			"PyYAML==6.0.1",
			"pyzmq==25.1.2",
			"qdrant-client==1.11.3",
			"rank-bm25==0.2.2",
			"rapidfuzz==3.6.1",
			"ratelimiter==1.2.0.post0",
			"redis==5.2.0",
			"referencing==0.33.0",
			"regex==2024.9.11",
			"requests==2.32.3",
			"requests-oauthlib==1.3.1",
			"requests-toolbelt==1.0.0",
			"responses==0.18.0",
			"retry==0.9.2",
			"rich==13.7.1",
			"rpds-py==0.18.0",
			"rq==2.0.0",
			"rsa==4.9",
			"ruff==0.1.15",
			"s3transfer==0.10.0",
			"sacrebleu==2.4.1",
			"sacremoses==0.1.1",
			"safetensors==0.4.2",
			"sagemaker==2.232.1",
			"sagemaker-core==1.0.9",
			"schema==0.7.7",
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
		okareoApiKey := os.Getenv("OKAREO_API_KEY")
		
		if okareoApiKey != "" {
			env = append(env, "OTEL_ENDPOINT=https://api.okareo.com/v0/traces")
			env = append(env, "OTEL_HEADERS=api-key="+okareoApiKey)
			env = append(env, "OTEL_EXPORTER=otlp_http")
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
	proxyCmd.Flags().StringP("config", "c", "", "Path to config file")
}