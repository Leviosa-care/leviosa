package config

import (
	"fmt"

	"github.com/hengadev/leviosa/pkg/envmode"
)

const BaseBucketName = "leviosa-assets"

type s3Creds struct {
	BucketName string
}

func (c *Config) GetS3() *s3Creds {
	return c.s3
}

func (c *Config) setS3(env envmode.Mode) error {
	switch env {
	case envmode.Dev, envmode.Staging:
		c.s3.BucketName = fmt.Sprintf("staging-%s", BaseBucketName)
	case envmode.Prod:
		c.s3.BucketName = fmt.Sprintf("production-%s", env.String(), BaseBucketName)
	default:
		return fmt.Errorf("mode value can only be 'development', 'production' or 'staging', got : %q", env)
	}
	return nil
}
