package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type ServiceKey struct {
	URL         string `json:"url"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
	Encryption  struct {
		ClientPrivateKey string `json:"client_private_key"`
	} `json:"encryption"`
}

type Parameters struct {
	Permission  int32
	Credentials []Credential
}

type Credential struct {
	Namespace string `yaml:"namespace,omitempty"`
	Type      string `yaml:"type,omitempty"`
	Name      string `yaml:"name,omitempty"`
	FileName  string `yaml:"fileName,omitempty"`
	Mode      *int32 `yaml:"mode,omitempty"`
}

func ParseServiceKey(jsonBytes []byte) (ServiceKey, error) {
	serviceKey := ServiceKey{}
	if err := json.Unmarshal(jsonBytes, &serviceKey); err != nil {
		return ServiceKey{}, fmt.Errorf("could not parse service key: %v", err)
	}

	return serviceKey, nil
}

func ParseParameters(attributes, permission string) (Parameters, error) {
	params := Parameters{}

	if err := json.Unmarshal([]byte(permission), &params.Permission); err != nil {
		return Parameters{}, fmt.Errorf("could not parse permission field: %v", err)
	}

	var err error
	params.Credentials, err = parseCredentials(attributes)
	if err != nil {
		return Parameters{}, err
	}

	if err = validateCredentials(params); err != nil {
		return Parameters{}, err
	}

	return params, nil
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

func validateCredentials(params Parameters) error {
	fileNames := make(map[string]bool)
	for _, cred := range params.Credentials {
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
