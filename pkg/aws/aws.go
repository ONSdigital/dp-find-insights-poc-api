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

// ProvideAWS is a wrapper around New that does not conk out
// if there is an error.
// We don't always need AWS support, so it's not necessarily
// an error if the AWS environment variable aren't set.
// Methods of *Clients must be prepared to see a nil
// receiver.
func ProvideAWS() *Clients {
	aws, _ := New()
	return aws
}
