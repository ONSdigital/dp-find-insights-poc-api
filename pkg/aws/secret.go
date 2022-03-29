package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/ONSdigital/log.go/v2/log"
)

// GetSecret looks up the secret named in arn.
// arn looks like: "arn:aws:secretsmanager:eu-central-1:352437599875:secret:fi-pg-x8rw4a"
func (clients *Clients) GetSecret(ctx context.Context, arn string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(arn),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := clients.sm.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Error(ctx, "AWS error", aerr, log.Data{"code": aerr.Code(), "message": aerr.Message()})
		} else {
			log.Error(ctx, "AWS error", err)
		}
		return "", err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return secretString, nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Error(ctx, "base64 decode", err)
			return "", err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		return decodedBinarySecret, nil
	}
}
