package main

import (
	"os"
	"testing"

	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) {
	assert.Equal(t, nil, os.Setenv("S3_ENDPOINT", "play.min.io"))
	assert.Equal(t, nil, os.Setenv("S3_ACCESS_KEY", "Q3AM3UQ867SPQQA43P2F"))
	assert.Equal(t, nil, os.Setenv("S3_SECRET_KEY", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"))

	s3, err := MinioFromEnv()
	assert.Equal(t, nil, err)

	if ok, err := s3.BucketExists("dcron-test-bucket"); !ok || err != nil {
		s3.MakeBucket("dcron-test-bucket", "")
	}

	n, err := s3.FPutObject("dcron-test-bucket", "hello_world.yaml", "./manifests/hello_world.yaml", minio.PutObjectOptions{})

	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, n)
}

func TestMinioFromEnv(t *testing.T) {
	setup(t)

	_, err := MinioFromEnv()

	assert.Equal(t, nil, err)
}

func TestReadJobSpecs(t *testing.T) {
	setup(t)

	jobs, err := ReadJobSpecs("dcron-test-bucket", "")

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(jobs))

}
