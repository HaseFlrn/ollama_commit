package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Model       string
	Ollama_Port int16
	Temperature float32
}

func GetConfig() Config {
	configPath := getConfigPath()

	if !configExists(configPath) {
		fmt.Println("Config not found. Writing defaults.")
		writeDefaultsConfig(configPath)
	}

	fmt.Println("Config found. Reading file now.")
	config := readConfig(configPath)

	return config
}

func getConfigPath() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get user home.")
		os.Exit(1)
	}

	configPath := userHome + "/.ollama_commit.conf"

	return configPath

}

func configExists(path string) bool {
	res := false
	if _, err := os.Stat(path); err == nil {
		res = true
	} else if errors.Is(err, os.ErrNotExist) {
		res = false
	} else {
		fmt.Printf("Something weird happened:\n%v", err)
		os.Exit(1)
	}
	return res
}

func readConfig(path string) Config {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read config file:\n%v", err)
	}

	config := Config{}

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Failed to read config file:\n%v", err)
	}

	fmt.Printf("%+v\n", config)
	return config
}

func writeDefaultsConfig(path string) {
	defaultConfig := Config{
		Model:       "llama3",
		Ollama_Port: 11434,
		Temperature: 1.0,
	}
	writeConfig(path, defaultConfig)
}

func writeConfig(path string, config Config) {

	configJson, _ := json.Marshal(config)

	os.WriteFile(path, configJson, 0644)
}

// Todo
func UpdateConfig(updateConfigDTO Config) {
	configPath := getConfigPath()
	updateConfig(configPath, updateConfigDTO)
}

func updateConfig(path string, updatedConfig Config) {}
