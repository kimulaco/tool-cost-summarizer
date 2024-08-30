package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type Config struct {
	Config   aws.Config
	IsLoaded bool
}

func NewConfig(accessKey string, secretAccessKey string) (Config, error) {
	if accessKey == "" || secretAccessKey == "" {
		return Config{IsLoaded: false}, errors.New("accessKey and secretAccessKey is required")
	}

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds))
	if err != nil {
		return Config{IsLoaded: false}, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return Config{Config: cfg, IsLoaded: true}, nil
}
