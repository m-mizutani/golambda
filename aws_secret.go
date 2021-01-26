package golambda

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return secretsmanager.New(ssn), nil
}

// SecretsManagerClient is wrapper of secretsmanager.SecretsManager
type SecretsManagerClient interface {
	GetSecretValue(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

// SecretsManagerFactory is factory function type to replace SecretsManagerClient
type SecretsManagerFactory func(region string) (SecretsManagerClient, error)

// GetSecretValuesWithFactory can call SecretsManager.GetSecretValue with your SecretsManagerClient by factory. It uses newDefaultSecretsManager if factory is nil
func GetSecretValuesWithFactory(secretArn string, values interface{}, factory SecretsManagerFactory) error {
	// sample: arn:aws:secretsmanager:ap-northeast-1:1234567890:secret:mytest
	arn := strings.Split(secretArn, ":")
	if len(arn) != 7 {
		return NewError("Invalid SecretsManager ARN format").With("arn", secretArn)
	}
	region := arn[3]

	if factory == nil {
		factory = newDefaultSecretsManager
	}
	mgr, err := factory(region)
	if err != nil {
		return WrapError(err).With("region", region)
	}

	result, err := mgr.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	})
	if err != nil {
		return WrapError(err, "Fail to retrieve secret values").With("arn", secretArn)
	}

	if err := json.Unmarshal([]byte(*result.SecretString), values); err != nil {
		return WrapError(err, "Fail to parse secret values as JSON").
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

// GetSecretValue is mock method for SecretsManagerMock. It checks if the secretId (ARN) exists in SecretsManagerMock.Secrets as key. It returns a string value if extsting or ResourceNotFoundException error if not existing.
func (x *SecretsManagerMock) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	x.Input = append(x.Input, input)
	value, ok := x.Secrets[*input.SecretId]
	if !ok {
		return nil, errors.New(secretsmanager.ErrCodeResourceNotFoundException)
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(value),
	}, nil
}

// NewSecretsManagerMock returns both of mock and factory method of the mock for testing. Developper can set secrets value as JSON to SecretsManagerMock.Secrets with key (secretes ARN). Also the mock stores Region that is extracted from secretArn and Input of secretsmanager.GetSecretValue when invoking GetSecretValuesWithFactory.
func NewSecretsManagerMock() (*SecretsManagerMock, SecretsManagerFactory) {
	mock := &SecretsManagerMock{
		Secrets: make(map[string]string),
	}
	return mock, func(region string) (SecretsManagerClient, error) {
		mock.Region = region
		return mock, nil
	}
}
