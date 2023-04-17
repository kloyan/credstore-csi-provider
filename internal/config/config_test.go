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
  mode: 0644
- name: myKey
  type: key
  namespace: dev
  fileName: key.der
  mode: 0400
- name: myKeyDec
  type: key
  namespace: dev
  fileName: key-dec.der
  mode: 420
- name: myKeyNoMode
  type: key
  namespace: dev
  fileName: key-no-mode.der
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
		attributes map[string]string
		expected   Parameters
	}{
		{
			name:       "valid credentials",
			permission: "420",
			attributes: map[string]string{"credentials": credentials},
			expected: Parameters{
				Permission: 420,
				Credentials: []Credential{
					{"dev", "password", "myPassword", "password.txt", modePtr(0644)},
					{"dev", "key", "myKey", "key.der", modePtr(0400)},
					{"dev", "key", "myKeyDec", "key-dec.der", modePtr(0644)},
					{"dev", "key", "myKeyNoMode", "key-no-mode.der", nil},
				},
			},
		},
		{
			name:       "no credentials",
			permission: "420",
			attributes: nil,
			expected: Parameters{
				Permission:  420,
				Credentials: nil,
			},
		},
	}

	for _, d := range data {
		jsonStr, err := json.Marshal(d.attributes)
		require.NoError(t, err, d.name)

		actual, err := ParseParameters(string(jsonStr), d.permission)
		require.NoError(t, err, d.name)
		require.Equal(t, d.expected, actual, d.name)
	}
}

func TestParse_Errors(t *testing.T) {
	data := []struct {
		name       string
		permission string
		attributes map[string]string
		errorMsg   string
	}{
		{
			name:       "missing name",
			permission: "420",
			attributes: map[string]string{"credentials": noNameCredential},
			errorMsg:   "credential name cannot be empty",
		},
		{
			name:       "missing namespace",
			permission: "420",
			attributes: map[string]string{"credentials": noNamespaceCredential},
			errorMsg:   "credential namespace cannot be empty",
		},
		{
			name:       "missing type",
			permission: "420",
			attributes: map[string]string{"credentials": noTypeCredential},
			errorMsg:   "credential type cannot be empty or invalid",
		},
		{
			name:       "invalid type",
			permission: "420",
			attributes: map[string]string{"credentials": invalidTypeCredential},
			errorMsg:   "credential type cannot be empty or invalid",
		},
		{
			name:       "missing file name",
			permission: "420",
			attributes: map[string]string{"credentials": noFileNameCredential},
			errorMsg:   "credential file name cannot be empty",
		},
		{
			name:       "duplicate file name",
			permission: "420",
			attributes: map[string]string{"credentials": duplicateFileNames},
			errorMsg:   "file name must be unique, password.txt is duplicated",
		},
	}

	for _, d := range data {
		jsonStr, err := json.Marshal(d.attributes)
		require.NoError(t, err, d.name)

		actual, err := ParseParameters(string(jsonStr), d.permission)
		require.EqualError(t, err, d.errorMsg, d.name)
		require.Equal(t, Parameters{}, actual, d.name)
	}
}

func modePtr(mode int32) *int32 {
	return &mode
}
