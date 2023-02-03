package ecs

import (
	"testing"

	awsecs "github.com/aws/aws-sdk-go/service/ecs"
)

// Test behavior for when there are no eni details
func TestDetermineENIDetails_Empty(t *testing.T) {
	taskResult := awsecs.Task{
		Attachments: make([]*awsecs.Attachment, 0),
	}

	found, _, _ := determineENIDetails(&taskResult)

	if found {
		t.Error("The blank attachment should not find an eni")
	}
}

// Test behavior for one simple eni attachment
func TestDetermineENIDetails_SimpleENI(t *testing.T) {
	taskResult := awsecs.Task{
		Attachments: make([]*awsecs.Attachment, 1),
	}

	taskResult.Attachments[0] = &awsecs.Attachment{
		Details: make([]*awsecs.KeyValuePair, 2),
	}

	var expectedSubnetName, expectedENIName = detailSubnetId, detailNetworkInterfaceId
	var expectedSubnet, expectedENIID = "abc", "123"
	var eniType = eniAttachmentType

	taskResult.Attachments[0].Type = &eniType
	taskResult.Attachments[0].Details[0] = &awsecs.KeyValuePair{
		Name:  &expectedSubnetName,
		Value: &expectedSubnet,
	}
	taskResult.Attachments[0].Details[1] = &awsecs.KeyValuePair{
		Name:  &expectedENIName,
		Value: &expectedENIID,
	}

	found, eniResult, subnetResult := determineENIDetails(&taskResult)

	if !found {
		t.Error("The blank attachment should find an eni")
	}

	if eniResult != expectedENIID {
		t.Errorf("Should find ENIID. Was %s expected %s", eniResult, expectedENIID)
	}

	if subnetResult != expectedSubnet {
		t.Errorf("Should find subnetid. Was %s expected %s", subnetResult, expectedSubnet)
	}
}

// Test behavior for when there are multiple attachments, the first being ENI, then Service Connect
func TestDetermineENIDetails_ServiceConnect(t *testing.T) {
	taskResult := awsecs.Task{
		Attachments: make([]*awsecs.Attachment, 2),
	}

	taskResult.Attachments[0] = &awsecs.Attachment{
		Details: make([]*awsecs.KeyValuePair, 2),
	}

	taskResult.Attachments[1] = &awsecs.Attachment{
		Details: make([]*awsecs.KeyValuePair, 2),
	}

	var expectedSubnetName, expectedENIName = detailSubnetId, detailNetworkInterfaceId
	var expectedSubnet, expectedENIID = "abc", "123"
	var notExpectedSubnet, noExpectedENIID = "xyz", "456"
	var eniType, serviceConnectType = eniAttachmentType, "Service Connect"

	taskResult.Attachments[0].Type = &eniType
	taskResult.Attachments[0].Details[0] = &awsecs.KeyValuePair{
		Name:  &expectedSubnetName,
		Value: &expectedSubnet,
	}
	taskResult.Attachments[0].Details[1] = &awsecs.KeyValuePair{
		Name:  &expectedENIName,
		Value: &expectedENIID,
	}

	taskResult.Attachments[1].Type = &serviceConnectType
	taskResult.Attachments[1].Details[0] = &awsecs.KeyValuePair{
		Name:  &expectedSubnetName,
		Value: &notExpectedSubnet,
	}
	taskResult.Attachments[1].Details[1] = &awsecs.KeyValuePair{
		Name:  &expectedENIName,
		Value: &noExpectedENIID,
	}

	found, eniResult, subnetResult := determineENIDetails(&taskResult)

	if !found {
		t.Error("The blank attachment should find an eni")
	}

	if eniResult != expectedENIID {
		t.Errorf("Should find ENIID. Was %s expected %s", eniResult, expectedENIID)
	}

	if subnetResult != expectedSubnet {
		t.Errorf("Should find subnetid. Was %s expected %s", subnetResult, expectedSubnet)
	}
}
