package golambda_test

import (
	"errors"
	"testing"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/golambda"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummySecret struct {
	Myth string `json:"myth"`
}

func TestGetSecretWithFactory(t *testing.T) {
	secretARN := "arn:aws:secretsmanager:us-east-0:1234567890:secret:mytest"
	mock, newMock := golambda.NewSecretsManagerMock()
	mock.Secrets[secretARN] = `{"myth":"magic"}`

	t.Run("can get secret values with custom SecretsManagerClient", func(t *testing.T) {
		var result dummySecret
		err := golambda.GetSecretValuesWithFactory(secretARN, &result, newMock)
		require.NoError(t, err)
		assert.Equal(t, "us-east-0", mock.Region)
		assert.Equal(t, "magic", result.Myth)
	})

	t.Run("fail when SecretsManagerFactory returns error", func(t *testing.T) {
		var result dummySecret
		newErr := goerr.New("something wrong")
		err := golambda.GetSecretValuesWithFactory(secretARN, &result, func(region string) (golambda.SecretsManagerClient, error) {
			return nil, newErr
		})

		require.Error(t, err)
		require.True(t, errors.Is(err, newErr))
	})
}
