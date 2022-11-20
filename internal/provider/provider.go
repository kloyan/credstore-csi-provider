package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/kloyan/credstore-csi-provider/internal/client"
	"github.com/kloyan/credstore-csi-provider/internal/config"
	pb "sigs.k8s.io/secrets-store-csi-driver/provider/v1alpha1"
)

type Provider struct {
	credStoreClient client.Client
}

func NewProvider(credStoreClient client.Client) *Provider {
	return &Provider{
		credStoreClient: credStoreClient,
	}
}

func (p *Provider) HandleMountRequest(ctx context.Context, params config.Parameters) (*pb.MountResponse, error) {
	var files []*pb.File
	var versions []*pb.ObjectVersion

	for _, cred := range params.Credentials {
		content, err := p.getCredentialContent(ctx, cred)
		if err != nil {
			return nil, err
		}

		files = append(files, &pb.File{
			Path:     cred.FileName,
			Mode:     int32(params.Permission),
			Contents: []byte(content),
		})

		version, err := generateVersion(cred, content)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	return &pb.MountResponse{
		Files:         files,
		ObjectVersion: versions,
	}, nil
}

func (p *Provider) getCredentialContent(ctx context.Context, cred config.Credential) (string, error) {
	if cred.Type == "password" {
		pass, err := p.credStoreClient.GetPassword(cred.Namespace, cred.Name)
		if err != nil {
			return "", err
		}
		return pass.Value, nil
	}

	if cred.Type == "key" {
		key, err := p.credStoreClient.GetKey(cred.Namespace, cred.Name)
		if err != nil {
			return "", err
		}
		return key.Value, nil
	}

	return "", fmt.Errorf("invalid credential type %s", cred.Type)
}

func generateVersion(cred config.Credential, content string) (*pb.ObjectVersion, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(fmt.Sprintf("%v:%s", cred, content)))
	if err != nil {
		return nil, err
	}

	return &pb.ObjectVersion{
		Id:      fmt.Sprintf("%s/%s/%s", cred.Namespace, cred.Type, cred.Name),
		Version: base64.URLEncoding.EncodeToString(hash.Sum(nil)),
	}, nil
}
