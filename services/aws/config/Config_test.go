package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name            string
	AccessKey       string
	SecretAccessKey string
	IsLoaded        bool
	Error           error
}

func TestNewConfig_Success(t *testing.T) {
	testCases := []testCase{
		{
			Name:            "Should content is correct",
			AccessKey:       "testAccessKey",
			SecretAccessKey: "testSecretKey",
			IsLoaded:        true,
			Error:           nil,
		},
		{
			Name:            "Should required accessKey",
			AccessKey:       "",
			SecretAccessKey: "testSecretKey",
			IsLoaded:        false,
			Error:           errors.New("accessKey and secretAccessKey is required"),
		},
		{
			Name:            "Should required secretAccessKey",
			AccessKey:       "testAccessKey",
			SecretAccessKey: "",
			IsLoaded:        false,
			Error:           errors.New("accessKey and secretAccessKey is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			config, err := NewConfig(tc.AccessKey, tc.SecretAccessKey)

			assert.Equal(t, err, tc.Error)
			assert.Equal(t, config.IsLoaded, tc.IsLoaded)
		})
	}
}
