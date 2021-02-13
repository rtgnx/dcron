package main

import (
	"os"
	"path"

	"github.com/minio/minio-go"
)

// MinioFromEnv returns minio client from environment variables
func MinioFromEnv() (*minio.Client, error) {
	return minio.New(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), true)
}

// ReadJobSpecs from S3
func ReadJobSpecs(bucket, prefix string) ([]JobSpec, error) {
	jobs := make([]JobSpec, 0)

	s3, _ := MinioFromEnv()

	if ok, err := s3.BucketExists(bucket); err != nil || !ok {
		return jobs, err
	}

	done := make(chan struct{})
	defer close(done)

	for object := range s3.ListObjects(bucket, prefix, false, done) {
		obj, _ := s3.GetObject(bucket, path.Join(s3Prefix, object.Key), minio.GetObjectOptions{})
		defer obj.Close()

		spec := new(JobSpec)
		if err := spec.fromReader(obj); err == nil {
			jobs = append(jobs, *spec)
		}
	}

	return jobs, nil
}
