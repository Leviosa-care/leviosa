package config

import (
	"github.com/hengadev/errsx"
)

type S3Secrets struct {
	BucketName string `json:"bucketname"`
}

func (s S3Secrets) Validate() error {
	var errs errsx.Map
	if s.BucketName == "" {
		errs.Set("S3 bucket name", "S3 bucket name cannot be empty")
	}
	return errs.AsError()
}

func (s S3Secrets) GetType() string {
	return "s3"
}
