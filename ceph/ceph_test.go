package ceph

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strings"
	"testing"
)

func TestPlayground(t *testing.T) {
	ctx := context.Background()
	container, err := RunContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	svc := CreateService(Config{
		AccessKey:  container.GetAccessKey(),
		SecretKey:  container.GetSecretKey(),
		Endpoint:   container.MustHttpURL(ctx),
		DisableSSL: true,
	})

	err = CreateBucket(svc, "my-bucket")
	if err != nil {
		t.Fatal(err)
	}

	err = CreateObject(svc, "my-bucket", "my-key", "Hello, World!")
	if err != nil {
		t.Fatal(err)
	}

	buckets, err := ListBuckets(svc)
	if err != nil {
		t.Fatal(err)
	}
	if len(buckets) != 2 {
		t.Fail()
	}

	content, err := GetObject(svc, "my-bucket", "my-key")
	if err != nil {
		t.Fatal(err)
	}
	if content != "Hello, World!" {
		t.FailNow()
	}
}

type Config struct {
	AccessKey  string
	SecretKey  string
	Endpoint   string
	DisableSSL bool
}

func CreateService(c Config) *s3.S3 {
	awsSession := session.Must(session.NewSession())
	awsCredentials := credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	awsConfig := aws.NewConfig().
		WithRegion("us-east-1").
		WithEndpoint(c.Endpoint).
		WithDisableSSL(c.DisableSSL).
		WithS3ForcePathStyle(true).
		WithCredentials(awsCredentials)

	return s3.New(awsSession, awsConfig)
}

func CreateBucket(svc *s3.S3, bucket string) error {
	input := &s3.CreateBucketInput{Bucket: &bucket}
	_, err := svc.CreateBucket(input)
	return err
}

func DeleteBucket(svc *s3.S3, bucket string) error {
	input := &s3.DeleteBucketInput{Bucket: &bucket}
	_, err := svc.DeleteBucket(input)
	return err
}

func ListBuckets(svc *s3.S3) ([]string, error) {
	var buckets []string
	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	for _, b := range result.Buckets {
		buckets = append(buckets, aws.StringValue(b.Name))
	}

	return buckets, err
}

func CreateObject(svc *s3.S3, bucket string, key string, content string) error {
	input := &s3.PutObjectInput{Body: strings.NewReader(content), Bucket: &bucket, Key: &key}
	_, err := svc.PutObject(input)
	return err
}

func GetObject(svc *s3.S3, bucket string, key string) (string, error) {
	input := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
	result, err := svc.GetObject(input)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, result.Body); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func DeleteObject(svc *s3.S3, bucket string, key string) error {
	input := &s3.DeleteObjectInput{Bucket: &bucket, Key: &key}
	_, err := svc.DeleteObject(input)
	return err
}
