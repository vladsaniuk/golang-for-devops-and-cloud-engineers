package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const (
	keyPairName           string = "ec2-key"
	ubuntuImageNameFilter string = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"
	canonicalsId          string = "099720109477"
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

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("Error constructing AWS config: " + err.Error())
	}
	ec2Client := ec2.NewFromConfig(cfg)

	ubuntuAmiId, err := getAmiId(ctx, ec2Client)
	if err != nil {
		slog.Error("Error getting list of image IDs by filter: " + err.Error())
	}

	keyPairFound, err := lookUpKeyPair(ctx, ec2Client)
	if err != nil {
		slog.Error("Error getting list of key pairs by filter: " + err.Error())
	}

	if !keyPairFound {
		keyPairCreatedOutput, err := createKeyPair(ctx, ec2Client)
		if err != nil {
			slog.Error("Error creating key pair: " + err.Error())
		}

		slog.Debug("Key pair created: " + *keyPairCreatedOutput.KeyName)
	}

	ec2RunOutput, err := createEc2Instance(ctx, ec2Client, ubuntuAmiId)
	if err != nil {
		slog.Error("Error starting EC2 instance: " + err.Error())
	}

	for _, ec2instance := range ec2RunOutput.Instances {
		slog.Debug("Instance started: " + *ec2instance.InstanceId)
	}
}
