package aws

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// GetSecret looks up the secret named in arn.
// arn looks like: "arn:aws:secretsmanager:eu-central-1:352437599875:secret:fi-pg-x8rw4a"
func (clients *Clients) GetSecret(arn string) (string, error) {
	if clients == nil {
		return "", errors.New("cannot get secret: AWS not initialised")
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(arn),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := clients.sm.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println(aerr.Error())
		} else {
			log.Println(err.Error())
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
			log.Println("Base64 Decode Error:", err)
			return "", err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		return decodedBinarySecret, nil
	}
}
