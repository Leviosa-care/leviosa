package config

import (
	"fmt"

	cfg "github.com/hengadev/leviosa/pkg/config"
	"github.com/hengadev/leviosa/pkg/envmode"
)

func (c *Config) GetS3() *cfg.S3Secrets {
	return c.s3
}

func (c *Config) setS3(env envmode.Mode) error {
	bucketname := c.viper.GetString("s3.bucketname")
	switch env {
	case envmode.Dev, envmode.Staging:
		c.s3.BucketName = fmt.Sprintf("staging-%s", bucketname)
	case envmode.Prod:
		c.s3.BucketName = fmt.Sprintf("production-%s", env.String(), bucketname)
	default:
		return fmt.Errorf("mode value can only be 'development', 'production' or 'staging', got : %q", env)
	}
	return nil
}
