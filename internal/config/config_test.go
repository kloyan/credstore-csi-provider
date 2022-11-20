package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	credentials = `
- name: myPassword
  type: password
  namespace: dev
  fileName: password.txt
- name: myKey
  type: key
  namespace: dev
  fileName: key.der
`
	noNameCredential = `
- type: password
  namespace: dev
  fileName: password.txt
`
	noNamespaceCredential = `
- name: myPassword
  type: password
  fileName: password.txt
`

	noTypeCredential = `
- name: myPassword
  namespace: dev
  fileName: password.txt
`

	invalidTypeCredential = `
- name: myPassword
  type: invalid
  namespace: dev
  fileName: password.txt
`

	noFileNameCredential = `
- name: myPassword
  type: password
  namespace: dev
`

	duplicateFileNames = `
- name: myPassword
  type: password
  namespace: dev
  fileName: password.txt
- name: myKey
  type: key
  namespace: dev
  fileName: password.txt
`
)

func TestParse(t *testing.T) {
	data := []struct {
		name       string
		permission string
		targetPath string
		attributes map[string]string
		expected   Config
	}{
		{
			name:       "valid credentials",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": credentials},
			expected: Config{
				Permission: 640,
				TargetPath: "my/path",
				Credentials: []Credential{
					{"dev", "password", "myPassword", "password.txt"},
					{"dev", "key", "myKey", "key.der"},
				},
			},
		},
		{
			name:       "no credentials",
			permission: "640",
			targetPath: "my/path",
			attributes: nil,
			expected: Config{
				Permission:  640,
				TargetPath:  "my/path",
				Credentials: nil,
			},
		},
	}

	for _, d := range data {
		jsonStr, err := json.Marshal(d.attributes)
		require.NoError(t, err, d.name)

		actual, err := Parse(string(jsonStr), d.targetPath, d.permission)
		require.NoError(t, err, d.name)
		require.Equal(t, d.expected, actual, d.name)
	}
}

func TestParse_Errors(t *testing.T) {
	data := []struct {
		name       string
		permission string
		targetPath string
		attributes map[string]string
		errorMsg   string
	}{
		{
			name:       "missing name",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": noNameCredential},
			errorMsg:   "credential name cannot be empty",
		},
		{
			name:       "missing namespace",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": noNamespaceCredential},
			errorMsg:   "credential namespace cannot be empty",
		},
		{
			name:       "missing type",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": noTypeCredential},
			errorMsg:   "credential type cannot be empty or invalid",
		},
		{
			name:       "invalid type",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": invalidTypeCredential},
			errorMsg:   "credential type cannot be empty or invalid",
		},
		{
			name:       "missing file name",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": noFileNameCredential},
			errorMsg:   "credential file name cannot be empty",
		},
		{
			name:       "duplicate file name",
			permission: "640",
			targetPath: "my/path",
			attributes: map[string]string{"credentials": duplicateFileNames},
			errorMsg:   "file name must be unique, password.txt is duplicated",
		},
		{
			name:       "missing target path",
			permission: "640",
			errorMsg:   "target path cannot be empty",
		},
	}

	for _, d := range data {
		jsonStr, err := json.Marshal(d.attributes)
		require.NoError(t, err, d.name)

		actual, err := Parse(string(jsonStr), d.targetPath, d.permission)
		require.EqualError(t, err, d.errorMsg, d.name)
		require.Equal(t, Config{}, actual, d.name)
	}
}
