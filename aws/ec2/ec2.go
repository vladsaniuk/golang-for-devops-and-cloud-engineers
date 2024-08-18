package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ec2Client interface {
	DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
	DescribeKeyPairs(ctx context.Context, params *ec2.DescribeKeyPairsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error)
	CreateKeyPair(ctx context.Context, params *ec2.CreateKeyPairInput, optFns ...func(*ec2.Options)) (*ec2.CreateKeyPairOutput, error)
	RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
}

func getAmiId(ctx context.Context, ec2Client ec2Client) (*string, error) {
	describeImagesOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Owners: []string{canonicalsId},
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{ubuntuImageNameFilter},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	ubuntuAmiId := describeImagesOutput.Images[0].ImageId
	return ubuntuAmiId, nil
}

func lookUpKeyPair(ctx context.Context, ec2Client ec2Client) (bool, error) {
	keyPairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("key-name"),
				Values: []string{keyPairName},
			},
		},
	})
	if err != nil {
		return false, err
	}

	keyPairFound := false
	for _, keyPair := range keyPairs.KeyPairs {
		if *keyPair.KeyName == keyPairName {
			keyPairFound = true
			break
		}
	}

	return keyPairFound, nil
}

func createKeyPair(ctx context.Context, ec2Client ec2Client) (*ec2.CreateKeyPairOutput, error) {
	keyPairCreatedOutput, err := ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyPairName),
	})
	if err != nil {
		return nil, err
	}

	return keyPairCreatedOutput, nil
}

func createEc2Instance(ctx context.Context, ec2Client ec2Client, ubuntuAmiId *string) (*ec2.RunInstancesOutput, error) {
	// run EC2 instance
	ec2RunOutput, err := ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
		ImageId:      ubuntuAmiId,
		InstanceType: types.InstanceTypeT3Micro,
		KeyName:      aws.String(keyPairName),
	})
	if err != nil {
		return nil, err
	}

	return ec2RunOutput, nil
}
