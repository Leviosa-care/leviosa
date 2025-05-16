package config

import (
	"fmt"

	mode "github.com/hengadev/leviosa/pkg/flags"
)

const BaseBucketName = "leviosa-assets"

type s3Creds struct {
	BucketName string
}

func (c *Config) GetS3() *s3Creds {
	return c.s3
}

func (c *Config) setS3(env mode.EnvMode) error {
	switch env {
	case mode.ModeDev, mode.ModeStaging:
		c.s3.BucketName = fmt.Sprintf("staging-%s", BaseBucketName)
	case mode.ModeProd:
		c.s3.BucketName = fmt.Sprintf("production-%s", env.String(), BaseBucketName)
	default:
		return fmt.Errorf("mode value can only be 'development', 'production' or 'staging', got : %q", env)
	}
	return nil
}
