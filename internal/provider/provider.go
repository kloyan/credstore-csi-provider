package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/kloyan/credstore-csi-provider/internal/client"
	"github.com/kloyan/credstore-csi-provider/internal/config"
	pb "sigs.k8s.io/secrets-store-csi-driver/provider/v1alpha1"
)

type Provider struct {
	credStoreClient *client.Client
}

func NewProvider(credStoreClient *client.Client) *Provider {
	return &Provider{
		credStoreClient: credStoreClient,
	}
}

func (p *Provider) HandleMountRequest(ctx context.Context, params config.Parameters) (*pb.MountResponse, error) {
	files := make([]*pb.File, len(params.Credentials))
	versions := make([]*pb.ObjectVersion, len(params.Credentials))
	errs := make([]error, len(params.Credentials))
	wg := sync.WaitGroup{}

	for i, cred := range params.Credentials {
		wg.Add(1)

		go func(i int, cred config.Credential) {
			defer wg.Done()

			content, err := p.getCredentialContent(ctx, cred)
			if err != nil {
				errs[i] = err
				return
			}

			mode := params.Permission
			if cred.Mode != nil {
				mode = *cred.Mode
			}

			files[i] = &pb.File{
				Path:     cred.FileName,
				Mode:     mode,
				Contents: []byte(content),
			}
			versions[i] = generateVersion(cred, content)
		}(i, cred)
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &pb.MountResponse{
		Files:         files,
		ObjectVersion: versions,
	}, nil
}

func (p *Provider) getCredentialContent(ctx context.Context, cred config.Credential) (string, error) {
	if cred.Type == "password" {
		pass, err := p.credStoreClient.GetPassword(ctx, cred.Namespace, cred.Name)
		if err != nil {
			return "", err
		}

		return pass.Value, nil
	}

	if cred.Type == "key" {
		key, err := p.credStoreClient.GetKey(ctx, cred.Namespace, cred.Name)
		if err != nil {
			return "", err
		}

		return key.Value, nil
	}

	return "", fmt.Errorf("invalid credential type %s", cred.Type)
}

func generateVersion(cred config.Credential, content string) *pb.ObjectVersion {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v:%s", cred, content)))

	return &pb.ObjectVersion{
		Id:      fmt.Sprintf("%s/%s/%s", cred.Namespace, cred.Type, cred.Name),
		Version: base64.URLEncoding.EncodeToString(hash.Sum(nil)),
	}
}
