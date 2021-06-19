package secretmanager

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManager struct {
	projectID string
	client    *secretmanager.Client
}

func (s *SecretManager) Close() {
	_ = s.client.Close()
}

func (s *SecretManager) GetSecret(ctx context.Context, key string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: "projects/" + s.projectID + "/secrets/" + key + "/versions/latest"}
	v, err := s.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("secretmanager: %v", err)
	}
	return string(v.GetPayload().GetData()), nil
}

func New(ctx context.Context, projectID string) (*SecretManager, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secretmanager: %v", err)
	}
	return &SecretManager{
		projectID: projectID,
		client:    c,
	}, nil
}
