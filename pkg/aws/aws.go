package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Clients struct {
	sess *session.Session
	sm   *secretsmanager.SecretsManager
}

func New() (*Clients, error) {
	cfg := aws.NewConfig()
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot get session: %w", err)
	}

	// Create Secrets Manager client
	sm := secretsmanager.New(sess, cfg)

	return &Clients{
		sess: sess,
		sm:   sm,
	}, nil
}
