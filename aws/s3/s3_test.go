package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
)

type mockS3Client struct {
	listBucketsOutput  *s3.ListBucketsOutput
	createBucketOutput *s3.CreateBucketOutput
	listObjectsOutput  *s3.ListObjectsOutput
}

func (m *mockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.listBucketsOutput, nil
}

func (m *mockS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.createBucketOutput, nil
}

func (m *mockS3Client) ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	return m.listObjectsOutput, nil
}

type mockS3Uploader struct {
	uploadOutput *manager.UploadOutput
}

func (m *mockS3Uploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	return m.uploadOutput, nil
}

type mockS3Downloader struct {
	numBytesDownloaded int64
}

func (m *mockS3Downloader) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error) {
	return m.numBytesDownloaded, nil
}

func TestLookupBucket(t *testing.T) {
	ctx := context.TODO()
	s3Client := &mockS3Client{
		listBucketsOutput: &s3.ListBucketsOutput{
			Buckets: []types.Bucket{
				{
					Name: aws.String(bucketName),
				},
			},
		},
	}

	bucketFound, err := lookupBucket(ctx, s3Client)
	if err != nil {
		t.Error("Error looking up bucket: " + err.Error())
	}

	if !bucketFound {
		t.Errorf("bucketFound is %v, expected to be true", bucketFound)
	}
}

func TestCreateBucket(t *testing.T) {
	ctx := context.TODO()
	s3Client := &mockS3Client{
		createBucketOutput: &s3.CreateBucketOutput{
			ResultMetadata: middleware.Metadata{},
		},
	}

	_, err := createBucket(ctx, s3Client)
	if err != nil {
		t.Error("Error creating bucket: " + err.Error())
	}
}

func TestLookupObject(t *testing.T) {
	var download string = "7abd75ad-42d3-446a-9852-b0dae9325bd7.txt"

	ctx := context.TODO()
	s3Client := &mockS3Client{
		listObjectsOutput: &s3.ListObjectsOutput{
			Contents: []types.Object{
				{
					Key: aws.String(fmt.Sprintf("%s/%s", uploadToDir, download)),
				},
			},
		},
	}

	objectFound, err := lookupObject(ctx, s3Client, download)
	if err != nil {
		t.Error("Error listing objects: " + err.Error())
	}

	if !objectFound {
		t.Errorf("Objects %s wasn't found in %s bucket, expected to be true", download, bucketName)
	}
}

func TestUploadFiles(t *testing.T) {
	var mockFile string = "75e21680-2ad7-4203-96bf-bf68906a088a.txt"
	err := os.WriteFile(fmt.Sprintf("%s/%s", filesDir, mockFile), []byte(mockFile), 0700)
	if err != nil {
		t.Error("Error creating test file: " + err.Error())
	}

	readFile, err := os.Open(fmt.Sprintf("%s/%s", filesDir, mockFile))
	if err != nil {
		t.Error("Error reading file " + mockFile + ": " + err.Error())
	}
	defer readFile.Close()

	ctx := context.TODO()
	s3Uploader := &mockS3Uploader{
		uploadOutput: &manager.UploadOutput{
			Key: aws.String(mockFile),
		},
	}
	uploadOutput, err := uploadFiles(ctx, s3Uploader, mockFile, readFile)
	if err != nil {
		t.Error("Error uploading file " + mockFile + ": " + err.Error())
	}

	if *uploadOutput.Key != mockFile {
		t.Errorf("uploadOutput.Key is %s, mockFile is %s, should be equal", *uploadOutput.Key, mockFile)
	}
}

func TestDownloadFile(t *testing.T) {
	var download string = "95124219-1ab2-4939-9bc1-cad22d413076.txt"

	ctx := context.TODO()
	s3Downloader := &mockS3Downloader{
		numBytesDownloaded: 1048576,
	}

	newFile, err := os.Create(fmt.Sprintf("downloads/%s", download))
	if err != nil {
		t.Error("Error creating new local file: " + err.Error())
	}
	defer newFile.Close()

	numBytesDownloaded, err := downloadFile(ctx, s3Downloader, newFile, download)
	if err != nil {
		t.Error("Error downloading file: " + err.Error())
	}

	if numBytesDownloaded == 0 {
		t.Error("numBytesDownloaded shouldn't be 0")
	}
}
