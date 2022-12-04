package client

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/kloyan/credstore-csi-provider/internal/config"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
)

type JWEDecryptor struct {
	privkey any
}

func NewJWEDecryptor(serviceKey config.ServiceKey) (JWEDecryptor, error) {
	bytes, err := base64.StdEncoding.DecodeString(serviceKey.Encryption.ClientPrivateKey)
	if err != nil {
		return JWEDecryptor{}, fmt.Errorf("could not decode private key: %v", err)
	}

	privkey, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		return JWEDecryptor{}, fmt.Errorf("could not parse private key: %v", err)
	}

	return JWEDecryptor{privkey: privkey}, nil
}

func (e *JWEDecryptor) Decrypt(data []byte) ([]byte, error) {
	decrypted, err := jwe.Decrypt(data, jwe.WithKey(jwa.RSA_OAEP_256, e.privkey))
	if err != nil {
		return nil, fmt.Errorf("could not decrypt data: %v", err)
	}

	return decrypted, nil
}
