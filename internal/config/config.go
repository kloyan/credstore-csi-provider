package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TargetPath  string
	Permission  os.FileMode
	Credentials []Credential
}

type Credential struct {
	Namespace string `yaml:"namespace,omitempty"`
	Type      string `yaml:"type,omitempty"`
	Name      string `yaml:"name,omitempty"`
	FileName  string `yaml:"fileName,omitempty"`
}

func Parse(attributes, targetPath, permission string) (Config, error) {
	config := Config{
		TargetPath: targetPath,
	}

	if err := json.Unmarshal([]byte(permission), &config.Permission); err != nil {
		return Config{}, fmt.Errorf("could not parse permission field: %v", err)
	}

	var err error
	config.Credentials, err = parseCredentials(attributes)
	if err != nil {
		return Config{}, err
	}

	if err = validate(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func parseCredentials(attributesStr string) ([]Credential, error) {
	var attributes map[string]string
	if err := json.Unmarshal([]byte(attributesStr), &attributes); err != nil {
		return nil, fmt.Errorf("could not parse attributes field: %v", err)
	}

	var creds []Credential
	credsYaml := attributes["credentials"]
	if err := yaml.Unmarshal([]byte(credsYaml), &creds); err != nil {
		return nil, fmt.Errorf("could not parse credentials field: %v", err)
	}

	return creds, nil
}

func validate(config Config) error {
	if len(config.TargetPath) == 0 {
		return fmt.Errorf("target path cannot be empty")
	}

	fileNames := make(map[string]bool)
	for _, cred := range config.Credentials {
		if len(cred.Namespace) == 0 {
			return fmt.Errorf("credential namespace cannot be empty")
		}

		if len(cred.Type) == 0 || (cred.Type != "password" && cred.Type != "key") {
			return fmt.Errorf("credential type cannot be empty or invalid")
		}

		if len(cred.Name) == 0 {
			return fmt.Errorf("credential name cannot be empty")
		}

		if len(cred.FileName) == 0 {
			return fmt.Errorf("credential file name cannot be empty")
		}

		if _, exists := fileNames[cred.FileName]; exists {
			return fmt.Errorf("file name must be unique, %s is duplicated", cred.FileName)
		}

		fileNames[cred.FileName] = true
	}

	return nil
}
