package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kloyan/credstore-csi-provider/internal/config"
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

type Client struct {
	BaseURL   string
	HTTP      *http.Client
	Decryptor JWEDecryptor
}

func NewClient(serviceKey config.ServiceKey, decryptor JWEDecryptor, timeout time.Duration) (*Client, error) {
	cert, err := tls.X509KeyPair([]byte(serviceKey.Certificate), []byte(serviceKey.Key))
	if err != nil {
		return nil, fmt.Errorf("could not parse x509 key pair: %v", err)
	}

	return &Client{
		BaseURL: serviceKey.URL,
		HTTP: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		},
		Decryptor: decryptor,
	}, nil
}

func (c *Client) GetPassword(ctx context.Context, namespace, name string) (*PasswordCredential, error) {
	url := fmt.Sprintf("%s/password?name=%s", c.BaseURL, name)
	password := &PasswordCredential{}
	err := c.getRequest(ctx, url, namespace, password)
	if err != nil {
		return nil, fmt.Errorf("could not get password %s/%s from credstore: %v", namespace, name, err)
	}

	return password, nil
}

func (c *Client) GetKey(ctx context.Context, namespace, name string) (*KeyCredential, error) {
	url := fmt.Sprintf("%s/key?name=%s", c.BaseURL, name)
	key := &KeyCredential{}
	err := c.getRequest(ctx, url, namespace, key)
	if err != nil {
		return nil, fmt.Errorf("could not get key %s/%s from credstore: %v", namespace, name, err)
	}

	return key, nil
}

func (c *Client) getRequest(ctx context.Context, url, namespace string, cred interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could not build http request: %v", err)
	}

	req.Header.Set("sapcp-credstore-namespace", namespace)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: got %v", resp.Status)
	}

	jwe, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %v", err)
	}

	decrypted, err := c.Decryptor.Decrypt(jwe)
	if err != nil {
		return fmt.Errorf("could not decrypt response body: %v", err)
	}

	err = json.Unmarshal(decrypted, cred)
	if err != nil {
		return fmt.Errorf("could not decode response body: %v", err)
	}

	return nil
}
