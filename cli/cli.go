package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type ConfigDetails struct {
	JenkinsURL string `json:"jenkins_url"`
	Username   string `json:"username"`
	Token      string `json:"token"`
}

var rootCmd = &cobra.Command{
	Use:   "jenkinsfile-validator",
	Short: "Validates Jenkinsfile",
	Long:  "Connects to Jenkinsfile Instance and Validates if the Jenkinsfile is valid",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Sets validator config",
	Long:  "Sets the information to validator to make sure to check if the Jenkinsfile is valid",
	Run:   configureSettings,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate File",
	Long:  "Validates the Jenkinsfile",
	Args:  cobra.MaximumNArgs(1),
	RunE:  validateFile,
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Jenkins connection",
	Long:  "Tests the connection to Jenkins using the configured credentials",
	RunE:  testConnection,
}

func testConnection(cmd *cobra.Command, args []string) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("Error loading config: %w", err)
	}

	if config.JenkinsURL == "" || config.Username == "" || config.Token == "" {
		return fmt.Errorf("Please configure Jenkins settings first using 'jenkinsfile-validator config'")
	}

	fmt.Println("Testing connection to Jenkins...")
	fmt.Printf("Jenkins URL: %s\n", config.JenkinsURL)
	fmt.Printf("Username: %s\n", config.Username)

	url := fmt.Sprintf("%s/api/json", config.JenkinsURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}
	req.SetBasicAuth(config.Username, config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error connecting to Jenkins: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("\nConnection successful!")
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			if nodeName, ok := result["nodeName"].(string); ok {
				fmt.Printf("Connected to Jenkins node: %s\n", nodeName)
			}
			if mode, ok := result["mode"].(string); ok {
				fmt.Printf("Jenkins mode: %s\n", mode)
			}
		}
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("Authentication failed. Please check your username and API token")
	case http.StatusForbidden:
		return fmt.Errorf("Access forbidden. Please check your permissions")
	default:
		return fmt.Errorf("Connection failed with status: %s", resp.Status)
	}
}

func validateFile(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		log.Fatal("Validate file command only takes in 1 args. Must pass path to Jenkinsfile in validator command.")
	}
	pathToJF := args[0]

	fileContent, err := os.ReadFile(pathToJF)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatal("Error when loading config file")
	}
	username := config.Username
	token := config.Token
	jenkinsurl := config.JenkinsURL
	endpoint := fmt.Sprintf("%s/pipeline-model-converter/validateJenkinsfile", jenkinsurl)

	// Prepare form data
	form := url.Values{}
	form.Set("jenkinsfile", string(fileContent))
	buf := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequest("POST", endpoint, buf)
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("Request failed with status: %s, response: %s", resp.Status, string(body[:n]))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("Error parsing response: %w", err)
	}

	status, ok := result["status"].(string)
	if !ok {
		return fmt.Errorf("Invalid response format: missing status field")
	}

	if status != "ok" {
		return fmt.Errorf("API request failed with status: %s", status)
	}

	// Check the actual validation result
	if data, ok := result["data"].(map[string]interface{}); ok {
		if resultStr, ok := data["result"].(string); ok {
			if resultStr == "success" {
				fmt.Println("✓ Jenkinsfile is valid")
				return nil
			} else if resultStr == "failure" {
				fmt.Println("✗ Jenkinsfile validation failed")

				if errors, ok := data["errors"].([]interface{}); ok && len(errors) > 0 {
					fmt.Println("\nErrors:")
					for _, err := range errors {
						if errMap, ok := err.(map[string]interface{}); ok {
							if errStr, ok := errMap["error"].(string); ok {
								fmt.Printf("- %s\n", errStr)
							}
						}
					}
				}

				return fmt.Errorf("Jenkinsfile validation failed")
			}
		}
	}

	return fmt.Errorf("Invalid response format: missing or invalid data field")
}

func init() {
	rootCmd.AddCommand(configCmd, validateCmd, testCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".validator_config.json"), nil
}

func loadConfig() (*ConfigDetails, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &ConfigDetails{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ConfigDetails
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *ConfigDetails) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func promptForInput(prompt, defaultValue string) string {
	fmt.Printf("%s", prompt)
	if defaultValue != "" {
		fmt.Printf(" [%s]", defaultValue)
	}
	fmt.Print(": ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}
	return input
}

func configureSettings(cmd *cobra.Command, args []string) {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("Configure Jenkins Validator Settings")
	fmt.Println("=====================================")

	if config.JenkinsURL != "" || config.Username != "" || config.Token != "" {
		fmt.Println("\nExisting configuration found:")
		fmt.Printf("Jenkins URL: %s\n", config.JenkinsURL)
		fmt.Printf("Username: %s\n", config.Username)
		if config.Token != "" {
			fmt.Printf("Token: %s\n", strings.Repeat("*", len(config.Token)))
		}
		fmt.Println()

		update := promptForInput("Do you want to update the configuration? (y/N)", "n")
		if strings.ToLower(update) != "y" && strings.ToLower(update) != "yes" {
			fmt.Println("Configuration unchanged.")
			return
		}
		fmt.Println()
	}

	config.JenkinsURL = promptForInput("Jenkins URL", config.JenkinsURL)
	config.Username = promptForInput("Username", config.Username)
	config.Token = promptForInput("API Token", "")

	if err := saveConfig(config); err != nil {
		log.Fatalf("Error saving config: %v", err)
	}

	fmt.Println("\nConfiguration saved successfully!")
}
