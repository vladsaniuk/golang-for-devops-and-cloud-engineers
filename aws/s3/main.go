package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	filesDir    string = "uploads"
	uploadToDir string = "reports"
	bucketName  string = "vladsanyuk-b961745f-fa36-4cfa-b9eb-7457e4c012fd"
)

func main() {
	// construct default logger
	var programLevel = new(slog.LevelVar) // Info by default
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(logger))

	// set log level to debug, if OS env DEBUG set as 1
	if os.Getenv("DEBUG") == "1" {
		programLevel.Set(slog.LevelDebug)
	}

	var (
		upload   bool
		download string
	)

	flag.BoolVar(&upload, "upload", false, "Bool, upload files from folder")
	flag.StringVar(&download, "download", "", "String, download specified file")
	flag.Parse()

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("Error constructing AWS config: " + err.Error())
	}
	s3Client := s3.NewFromConfig(cfg)

	if upload {
		bucketFound, err := lookupBucket(ctx, s3Client)
		if err != nil {
			slog.Error("Error looking up bucket: " + err.Error())
		}

		if !bucketFound {
			bucketCreatedOutput, err := createBucket(ctx, s3Client)
			if err != nil {
				slog.Error("Error creating bucket: " + err.Error())
			}

			slog.Debug("bucketCreatedOutput is: " + fmt.Sprintf("%v", bucketCreatedOutput.ResultMetadata))
		}

		s3Uploader := manager.NewUploader(s3Client)

		files := make([]string, 0)
		_ = filepath.WalkDir(filesDir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				slog.Debug("File is: " + d.Name())
				files = append(files, d.Name())
			}

			return nil
		})

		for index, file := range files {
			indexStr := strconv.Itoa(index)
			slog.Debug("file " + indexStr + " is: " + file)

			readFile, err := os.Open(fmt.Sprintf("%s/%s", filesDir, file))
			if err != nil {
				slog.Error("Error reading file " + file + ": " + err.Error())
			}
			defer readFile.Close()

			uploadOutput, err := uploadFiles(ctx, s3Uploader, file, readFile)
			if err != nil {
				slog.Error("Error uploading file " + file + ": " + err.Error())
			}

			slog.Debug(file + " file uploaded as: " + *uploadOutput.Key)
		}
	}

	if download != "" {
		objectFound, err := lookupObject(ctx, s3Client, download)
		if err != nil {
			slog.Error("Error listing objects: " + err.Error())
		}

		if !objectFound {
			slog.Error("Objects " + download + " wasn't found in " + bucketName + " bucket")
			os.Exit(1)
		}

		s3Downloader := manager.NewDownloader(s3Client)

		newFile, err := os.Create(fmt.Sprintf("downloads/%s", download))
		if err != nil {
			slog.Error("Error creating new local file: " + err.Error())
		}
		defer newFile.Close()

		numBytesDownloaded, err := downloadFile(ctx, s3Downloader, newFile, download)
		if err != nil {
			slog.Error("Error downloading file: " + err.Error())
		}

		if numBytesDownloaded != 0 {
			slog.Info("Successfully downloaded " + download)
		}
	}
}
