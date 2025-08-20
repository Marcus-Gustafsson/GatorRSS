// Package config provides functionality for reading and writing the application's
// configuration file, including database credentials and the current user. If the
// configuration file does not exist, a new file with default values will be
// created.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config represents the application's configuration file structure.
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// SetUser updates CurrentUserName and writes the new config to disk.
// Returns an error if marshalling or writing fails.
func (cfgPtr *Config) SetUser(userName string) error {

	if len(userName) == 0 {
		return errors.New("error: userName must be atleast 1 char")
	}

	cfgPtr.CurrentUserName = userName

	jsonData, err := json.MarshalIndent(cfgPtr, "", "  ")
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// getConfigFilePath returns the absolute path to the configuration file
// in the user's home directory.
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homePath, configFileName)
	return configPath, nil
}

// Read loads the configuration from the JSON file in the user's home directory.
// If the file does not exist, it creates one with default values and returns
// the default Config. Returns an error if any file operations or unmarshalling
// fail.
func Read() (Config, error) {
	jsonConfigPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(jsonConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := Config{}
			jsonData, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return Config{}, err
			}
			if err := os.WriteFile(jsonConfigPath, jsonData, 0644); err != nil {
				return Config{}, err
			}
			fmt.Printf("Created new config file at %s\n", jsonConfigPath)
			return cfg, nil
		}
		return Config{}, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(byteValue, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
