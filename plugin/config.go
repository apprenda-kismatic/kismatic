package plugin

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type BackendAuth struct {
	AuthToken string
}

type KismaticConfig struct {
	Auths map[string]BackendAuth
}

const configFileName = "config.json"

func readConfig() (*KismaticConfig, error) {

	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			// Create directory if it does not exist
			if err = os.Mkdir(configDir, 0600); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	configFile := filepath.Join(configDir, configFileName)
	if _, err = os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			// Config file does not exist, return empty config struct
			return &KismaticConfig{Auths: map[string]BackendAuth{}}, nil
		}
		return nil, err
	}

	configJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := &KismaticConfig{}
	err = json.Unmarshal(configJSON, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func writeConfig(config *KismaticConfig) error {
	configJSON, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(configDir, configFileName)

	return ioutil.WriteFile(configFile, configJSON, 0600)
}

func deleteConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configFile := filepath.Join(configDir, configFileName)

	return os.Remove(configFile)
}

func getConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(usr.HomeDir, ".kismatic")
	return configDir, nil
}
