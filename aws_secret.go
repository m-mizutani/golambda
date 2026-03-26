package golambda

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/m-mizutani/goerr"
)

// GetSecretValues bind secret data of AWS Secrets Manager to values. values should be set as pointer of struct with json meta tag.
//
//     type mySecret struct {
//         Token string `json:"token"`
//     }
//     var secret mySecret
//     if err := golambda.GetSecretValues(secretARN, &secret); err != nil {
//         log.Fatal("Failed: ", err)
//     }
func GetSecretValues(secretArn string, values interface{}) error {
	return GetSecretValuesWithFactory(secretArn, values, nil)
}

func newDefaultSecretsManager(region string) (SecretsManagerClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return secretsmanager.NewFromConfig(cfg), nil
}

// SecretsManagerClient is wrapper of secretsmanager.SecretsManager
type SecretsManagerClient interface {
	GetSecretValue(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretsManagerFactory is factory function type to replace SecretsManagerClient
type SecretsManagerFactory func(region string) (SecretsManagerClient, error)

// GetSecretValuesWithFactory can call SecretsManager.GetSecretValue with your SecretsManagerClient by factory. It uses newDefaultSecretsManager if factory is nil
func GetSecretValuesWithFactory(secretArn string, values interface{}, factory SecretsManagerFactory) error {
	// sample: arn:aws:secretsmanager:ap-northeast-1:1234567890:secret:mytest
	arn := strings.Split(secretArn, ":")
	if len(arn) != 7 {
		return goerr.Wrap(ErrInvalidARN).With("arn", secretArn)
	}
	region := arn[3]

	if factory == nil {
		factory = newDefaultSecretsManager
	}
	mgr, err := factory(region)
	if err != nil {
		return goerr.Wrap(err).With("region", region)
	}

	result, err := mgr.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	})
	if err != nil {
		return ErrFailedSecretsManager.Wrap(err).With("arn", secretArn)
	}

	if err := json.Unmarshal([]byte(aws.ToString(result.SecretString)), values); err != nil {
		return ErrFailedDecodeSecret.Wrap(err).
			With("arn", secretArn).
			With("GetSecretValue:result", result)
	}

	return nil
}

// SecretsManagerMock is mock of SecretsManagerClient for testing.
type SecretsManagerMock struct {
	Secrets map[string]string
	Region  string
	Input   []*secretsmanager.GetSecretValueInput
}

// GetSecretValue is mock method for SecretsManagerMock. It checks if the secretId (ARN) exists in SecretsManagerMock.Secrets as key. It returns a string value if existing or ResourceNotFoundException error if not existing.
func (x *SecretsManagerMock) GetSecretValue(_ context.Context, input *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	x.Input = append(x.Input, input)
	secretID := aws.ToString(input.SecretId)
	value, ok := x.Secrets[secretID]
	if !ok {
		return nil, errors.New("ResourceNotFoundException")
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(value),
	}, nil
}

// NewSecretsManagerMock returns both of mock and factory method of the mock for testing. Developer can set secrets value as JSON to SecretsManagerMock.Secrets with key (secretes ARN). Also the mock stores Region that is extracted from secretArn and Input of secretsmanager.GetSecretValue when invoking GetSecretValuesWithFactory.
func NewSecretsManagerMock() (*SecretsManagerMock, SecretsManagerFactory) {
	mock := &SecretsManagerMock{
		Secrets: make(map[string]string),
	}
	return mock, func(region string) (SecretsManagerClient, error) {
		mock.Region = region
		return mock, nil
	}
}
