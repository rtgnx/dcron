package main

import (
	"os"

	"github.com/minio/minio-go"
)

// MinioFromEnv returns minio client from environment variables
func MinioFromEnv() (*minio.Client, error) {
	return minio.New(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), true)
}
