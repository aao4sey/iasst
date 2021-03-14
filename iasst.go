package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "sg",
				Usage: "Shows list of resources to which specified security group are attached",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "security-group-id",
						Aliases:  []string{"id"},
						Usage:    "Sets security group id",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "check-security-group",
						Aliases: []string{"s"},
						Usage:   "Shows security groups to which the security group is being used as rule.",
					},
					&cli.BoolFlag{
						Name:    "check-eni",
						Aliases: []string{"e"},
						Usage:   "Shows ENIs to which the security group is being attached.",
					},
				},
				Action: checkSecurityGroupDependency,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func checkSecurityGroupDependency(c *cli.Context) error {
	targetSecurityGroupId := c.String("security-group-id")
	fmt.Printf("check target: %s\n\n", targetSecurityGroupId)

	if c.Bool("check-security-group") {
		checkUsedByOtherSecurityGroup(targetSecurityGroupId)
	}
	if c.Bool("check-eni") {
		checkUsedByEni(targetSecurityGroupId)
	}
	return nil
}

func checkUsedByOtherSecurityGroup(securityGroupId string) {
	securityGroups := getSecurityGroupList()

	fmt.Printf("TargetGroupId,GroupId,GroupName,FromPort,ToPort,IpProtocol\n")
	for _, sg := range securityGroups {

		extractedRules := extractReleventRulesById(sg.IpPermissions, securityGroupId)
		if len(extractedRules) == 0 {
			continue
		}
		for _, rule := range extractedRules {
			fmt.Printf("%s,%s,%s,%s,%s,%s\n",
				securityGroupId,
				*sg.GroupId,
				*sg.GroupName,
				rule.FromPort,
				rule.ToPort,
				rule.IpProtocol,
			)
		}
	}
}

func hasSecurityGroupId(groupIds []*ec2.GroupIdentifier, securityGroupId string) bool {
	for _, groupSet := range groupIds {
		if securityGroupId == *groupSet.GroupId {
			return true
		}
	}
	return false
}

func checkUsedByEni(securityGroupId string) {
	eniList := getEniList()
	eniListSecurityGroupAttached := []ec2.NetworkInterface{}
	for _, eni := range eniList {
		if hasSecurityGroupId(eni.Groups, securityGroupId) {
			eniListSecurityGroupAttached = append(eniListSecurityGroupAttached, eni)
		}
	}

	fmt.Printf("NetworkInterfaceId,Status,InterfaceType,InstanceId,InstanceOwnerId\n")
	for _, eni := range eniListSecurityGroupAttached {
		var instanceId string
		var instanceOwnerId string
		var interfaceType string

		if *eni.Status != "available" {
			if eni.Attachment.InstanceId == nil {
				instanceId = "None"
			} else {
				instanceId = *eni.Attachment.InstanceId
			}
			if eni.Attachment.InstanceOwnerId == nil {
				instanceOwnerId = "None"
			} else {
				instanceOwnerId = *eni.Attachment.InstanceOwnerId
			}
		} else {
			instanceId = "None"
			instanceOwnerId = "None"
		}

		if eni.InterfaceType == nil {
			interfaceType = "None"
		} else {
			interfaceType = *eni.InterfaceType
		}

		fmt.Printf("%s,%s,%s,%s,%s\n",
			*eni.NetworkInterfaceId,
			*eni.Status,
			interfaceType,
			instanceId,
			instanceOwnerId,
		)
	}
}

func extractReleventRulesById(ipPermissions []*ec2.IpPermission, securityGroupId string) []ResultRule {
	result := []ResultRule{}
	for _, rule := range ipPermissions {
		for _, userIdGroupPair := range rule.UserIdGroupPairs {
			if *userIdGroupPair.GroupId == securityGroupId {
				var fromPort string
				var toPort string
				var groupName string

				if rule.FromPort == nil {
					fromPort = "All"
				} else {
					fromPort = strconv.FormatInt(*rule.FromPort, 10)
				}

				if rule.ToPort == nil {
					toPort = "All"
				} else {
					toPort = strconv.FormatInt(*rule.ToPort, 10)
				}

				if userIdGroupPair.GroupName == nil {
					groupName = "None"
				} else {
					groupName = *userIdGroupPair.GroupName
				}

				result = append(result, ResultRule{
					FromPort:   fromPort,
					ToPort:     toPort,
					IpProtocol: *rule.IpProtocol,
					GroupId:    *userIdGroupPair.GroupId,
					GroupName:  groupName,
				})
			}
		}
	}
	return result
}

type ResultRule struct {
	FromPort   string
	ToPort     string
	IpProtocol string
	GroupId    string
	GroupName  string
}
