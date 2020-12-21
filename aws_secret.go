package golambda

import (
	"encoding/json"
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
	// sample: arn:aws:secretsmanager:ap-northeast-1:1234567890:secret:mytest
	arn := strings.Split(secretArn, ":")
	if len(arn) != 7 {
		return NewError("Invalid SecretsManager ARN format").With("arn", secretArn)
	}
	region := arn[3]

	ssn := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	mgr := secretsmanager.New(ssn)

	result, err := mgr.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	})
	if err != nil {
		return WrapError(err, "Fail to retrieve secret values").With("arn", secretArn)
	}

	if err := json.Unmarshal([]byte(*result.SecretString), &values); err != nil {
		return WrapError(err, "Fail to parse secret values as JSON").
			With("arn", secretArn).
			With("GetSecretValue:result", result)
	}

	return nil
}
