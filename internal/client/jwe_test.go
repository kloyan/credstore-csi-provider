package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"os"
	"testing"

	"github.com/kloyan/credstore-csi-provider/internal/config"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/stretchr/testify/require"
)

var (
	pubkey    any
	decryptor JWEDecryptor
	payload   = []byte("Hello, World!")
)

func TestMain(m *testing.M) {
	pubkeyBytes, err := os.ReadFile("mock/pubkey")
	if err != nil {
		panic(err)
	}

	pubkey, err = x509.ParsePKIXPublicKey(pubkeyBytes)
	if err != nil {
		panic(err)
	}

	privkeyBytes, err := os.ReadFile("mock/privkey")
	if err != nil {
		panic(err)
	}

	serviceKey := config.ServiceKey{
		Encryption: struct {
			ClientPrivateKey string "json:\"client_private_key\""
		}{
			ClientPrivateKey: base64.StdEncoding.EncodeToString(privkeyBytes),
		},
	}

	decryptor, err = NewJWEDecryptor(serviceKey)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestNewJWEDecryptor_InvalidKeys(t *testing.T) {
	data := []struct {
		name     string
		privkey  []byte
		errorMsg string
	}{
		{
			name:     "nil private key",
			privkey:  []byte(nil),
			errorMsg: "could not parse private key",
		},
		{
			name:     "invalid private key",
			privkey:  []byte("foobar"),
			errorMsg: "could not parse private key",
		},
		{
			name:     "not pkcs8 private key",
			privkey:  generatePKCS1Key(),
			errorMsg: "could not parse private key",
		},
	}

	for _, d := range data {
		serviceKey := config.ServiceKey{
			Encryption: struct {
				ClientPrivateKey string "json:\"client_private_key\""
			}{
				ClientPrivateKey: base64.StdEncoding.EncodeToString(d.privkey),
			},
		}

		decryptor, err := NewJWEDecryptor(serviceKey)
		require.Equal(t, JWEDecryptor{}, decryptor)
		require.ErrorContains(t, err, d.errorMsg)
	}
}

func TestDecrypt(t *testing.T) {
	encrypted, err := jwe.Encrypt(payload, jwe.WithJSON(), jwe.WithKey(jwa.RSA_OAEP_256, pubkey))
	require.NoError(t, err)
	require.NotEqual(t, payload, encrypted)

	decrypted, err := decryptor.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, payload, decrypted)
}

func TestDecrypt_InvalidData(t *testing.T) {
	invalidDecryptedPayload := []byte("foobar")

	decrypted, err := decryptor.Decrypt(invalidDecryptedPayload)
	require.ErrorContains(t, err, "could not decrypt data")
	require.Nil(t, decrypted)
}

func generatePKCS1Key() []byte {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	privkey := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(key))
	return []byte(privkey)
}
