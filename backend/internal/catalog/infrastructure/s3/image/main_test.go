package imageMedia_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/s3/image"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
)

var (
	pgContainer *tu.PostgresContainer
	s3Client    *s3.Client
	repo        ports.ImageMedia
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error

	// Postgres container
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup postgres container: %v", err))
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	log.Println("Setting S3 testcontainer...")
	// s3 container
	localstackContainer, err := tu.SetupLocalstack(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup S3 container: %v", err)
	}
	defer tu.TeardownLocalstack(ctx, nil, localstackContainer)
	log.Println("S3 testcontainer et.")

	// s3 config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"test", // Access Key ID
			"test", // Secret Access Key
			"",     // Session Token (empty for Localstack)
		)),
	)
	if err != nil {
		log.Fatalf("Load default S3 configuration: %s\n", err)
	}
	// s3 client
	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(localstackContainer.S3Endpoint)
		o.UsePathStyle = true // Required for Localstack
		o.Region = "us-east-1"
	})
	// Create a test bucket in Localstack S3
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(td.BUCKETNAME), // Use your test bucket name
	})
	if err != nil {
		log.Fatalf("Failed to create test S3 bucket: %v", err)
	}
	log.Println("Test S3 bucket created.")

	repo = imageMedia.New(ctx, s3Client, td.BUCKETNAME)

	// Run tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code) // Commented out to allow cleanup before exiting in some environments
	// The `m.Run()` call handles exiting with the correct code in `go test`
}
