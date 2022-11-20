package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kloyan/credstore-csi-provider/internal/config"
)

const (
	namespaceHeader = "sapcp-credstore-namespace"
)

type PasswordCredential struct {
	ID         string `json:"id"`
	ModifiedAt string `json:"modifiedAt"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Metadata   string `json:"metadata"`
	Value      string `json:"value"`
}

type KeyCredential struct {
	ID         string `json:"id"`
	ModifiedAt string `json:"modifiedAt"`
	Name       string `json:"name"`
	Format     string `json:"format"`
	Username   string `json:"username"`
	Metadata   string `json:"metadata"`
	Value      string `json:"value"`
}

type Client interface {
	GetPassword(namespace, name string) (*PasswordCredential, error)
	GetKey(namespace, name string) (*KeyCredential, error)
}

type client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(serviceKey config.ServiceKey, timeout time.Duration) (*client, error) {
	cert, err := tls.X509KeyPair([]byte(serviceKey.Certificate), []byte(serviceKey.Key))
	if err != nil {
		return nil, fmt.Errorf("could not parse x509 key pair: %v", err)
	}

	return &client{
		BaseURL: serviceKey.URL,
		HTTP: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		},
	}, nil
}

func (c *client) GetPassword(namespace, name string) (*PasswordCredential, error) {
	url := fmt.Sprintf("%s/password?name=%s", c.BaseURL, name)
	cred := &PasswordCredential{}
	err := c.getRequest(url, namespace, cred)
	return cred, err
}

func (c *client) GetKey(namespace, name string) (*KeyCredential, error) {
	url := fmt.Sprintf("%s/key?name=%s", c.BaseURL, name)
	cred := &KeyCredential{}
	err := c.getRequest(url, namespace, cred)
	return cred, err
}

func (c *client) getRequest(url, namespace string, cred interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could not build http request: %v", err)
	}

	req.Header.Set(namespaceHeader, namespace)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: got %v", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(cred)
	if err != nil {
		return fmt.Errorf("could not decode response body: %v", err)
	}

	return nil
}
