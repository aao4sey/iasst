package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getEniList() []ec2.NetworkInterface {
	svc := getEc2Client()
	eniList := []ec2.NetworkInterface{}
	nextToken := ""
	var input ec2.DescribeNetworkInterfacesInput
	for i := 0; ; i++ {
		if i == 0 {
			input = ec2.DescribeNetworkInterfacesInput{}
		} else {
			input = ec2.DescribeNetworkInterfacesInput{
				NextToken: &nextToken,
			}
		}
		result, err := svc.DescribeNetworkInterfaces(
			&input,
		)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
		}
		for _, eni := range result.NetworkInterfaces {
			eniList = append(eniList, *eni)
		}
		if result.NextToken == nil {
			break
		}
		nextToken = *result.NextToken
	}
	return eniList
}
