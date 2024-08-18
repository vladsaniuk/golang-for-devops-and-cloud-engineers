package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Client interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
}

type s3Uploader interface {
	Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

type s3Downloader interface {
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error)
}

func lookupBucket(ctx context.Context, s3Client s3Client) (bool, error) {
	bucketsList, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return false, err
	}

	bucketFound := false
	for _, bucket := range bucketsList.Buckets {
		if *bucket.Name == bucketName {
			bucketFound = true
			break
		}
	}

	return bucketFound, nil
}

func createBucket(ctx context.Context, s3Client s3Client) (*s3.CreateBucketOutput, error) {
	bucketCreatedOutput, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	return bucketCreatedOutput, nil
}

func lookupObject(ctx context.Context, s3Client s3Client, download string) (bool, error) {
	objectsList, err := s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(uploadToDir),
	})
	if err != nil {
		return false, err
	}

	objectFound := false
	for _, object := range objectsList.Contents {
		if *object.Key == fmt.Sprintf("%s/%s", uploadToDir, download) {
			objectFound = true
			break
		}
	}

	return objectFound, nil
}

func uploadFiles(ctx context.Context, s3Uploader s3Uploader, file string, readFile *os.File) (*manager.UploadOutput, error) {
	uploadOutput, err := s3Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", uploadToDir, file)),
		Body:   readFile,
	})
	if err != nil {
		return nil, err
	}

	return uploadOutput, nil
}

func downloadFile(ctx context.Context, s3Downloader s3Downloader, newFile *os.File, download string) (int64, error) {
	numBytesDownloaded, err := s3Downloader.Download(ctx, newFile, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", uploadToDir, download)),
	})
	if err != nil {
		return 0, err
	}

	return numBytesDownloaded, nil
}
