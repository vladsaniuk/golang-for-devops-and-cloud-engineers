package main

import (
	"context"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type mockEc2Client struct {
	describeImagesOutput   *ec2.DescribeImagesOutput
	describeKeyPairsOutput *ec2.DescribeKeyPairsOutput
	createKeyPairOutput    *ec2.CreateKeyPairOutput
	runInstancesOutput     *ec2.RunInstancesOutput
}

const mockImageId string = "prod-x7h6cigkuiul6"

func (m *mockEc2Client) DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error) {
	return m.describeImagesOutput, nil
}

func (m *mockEc2Client) DescribeKeyPairs(ctx context.Context, params *ec2.DescribeKeyPairsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error) {
	return m.describeKeyPairsOutput, nil
}

func (m *mockEc2Client) CreateKeyPair(ctx context.Context, params *ec2.CreateKeyPairInput, optFns ...func(*ec2.Options)) (*ec2.CreateKeyPairOutput, error) {
	return m.createKeyPairOutput, nil
}

func (m *mockEc2Client) RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {
	return m.runInstancesOutput, nil
}

func TestGetAmiId(t *testing.T) {
	var mockImageId string = "prod-x7h6cigkuiul6"

	ctx := context.TODO()
	ec2Client := &mockEc2Client{
		describeImagesOutput: &ec2.DescribeImagesOutput{
			Images: []types.Image{
				{
					ImageId: aws.String(mockImageId),
				},
			},
		},
	}

	ubuntuAmiId, err := getAmiId(ctx, ec2Client)
	if err != nil {
		t.Error("Error getting list of image IDs by filter: " + err.Error())
	}

	if *ubuntuAmiId != mockImageId {
		t.Errorf("Image ID isn't correct, %s != %s", *ubuntuAmiId, mockImageId)
	}
}

func TestLookUpKeyPair(t *testing.T) {
	ctx := context.TODO()
	ec2Client := &mockEc2Client{
		describeKeyPairsOutput: &ec2.DescribeKeyPairsOutput{
			KeyPairs: []types.KeyPairInfo{
				{
					KeyName: aws.String(keyPairName),
				},
			},
		},
	}

	keyPairFound, err := lookUpKeyPair(ctx, ec2Client)
	if err != nil {
		t.Error("Error getting list of key pairs by filter: " + err.Error())
	}

	if !keyPairFound {
		t.Errorf("keyPairFound returned %v, expected true", keyPairFound)
	}
}

func TestCreateKeyPair(t *testing.T) {
	ctx := context.TODO()
	ec2Client := &mockEc2Client{
		createKeyPairOutput: &ec2.CreateKeyPairOutput{
			KeyName: aws.String(keyPairName),
		},
	}

	keyPairCreatedOutput, err := createKeyPair(ctx, ec2Client)
	if err != nil {
		t.Error("Error creating key pair: " + err.Error())
	}

	if *keyPairCreatedOutput.KeyName != keyPairName {
		t.Errorf("keyPairCreatedOutput.KeyName is %s, while expected is %s", *keyPairCreatedOutput.KeyName, keyPairName)
	}
}

func TestCreateEc2Instance(t *testing.T) {
	var mockInstanceId string = "i-0f3f71c5c31adaae2"
	ctx := context.TODO()
	ec2Client := &mockEc2Client{
		runInstancesOutput: &ec2.RunInstancesOutput{
			Instances: []types.Instance{
				{
					InstanceId: aws.String(mockInstanceId),
				},
			},
		},
	}
	ec2RunOutput, err := createEc2Instance(ctx, ec2Client, aws.String(mockImageId))
	if err != nil {
		slog.Error("Error starting EC2 instance: " + err.Error())
	}

	for _, ec2instance := range ec2RunOutput.Instances {
		slog.Debug("Instance started: " + *ec2instance.InstanceId)
	}
}
