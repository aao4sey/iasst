package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func getEc2Client() *ec2.EC2 {
	region := "ap-northeast-1"
	s, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ec2.New(s, &aws.Config{Region: &region})
}

func getSecurityGroupList() []*ec2.SecurityGroup {
	svc := getEc2Client()
	result, err := svc.DescribeSecurityGroups(
		&ec2.DescribeSecurityGroupsInput{},
	)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroupId.Malformed":
				fallthrough
			case "InvalidGroup.NotFound":
				exitErrorf("%s.", aerr.Message())
			}
		}
		exitErrorf("Unable to get descriptions for security groups, %v", err)
	}
	return result.SecurityGroups
}
