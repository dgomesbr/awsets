package lister

import (
	"fmt"

	"github.com/trek10inc/awsets/context"

	"github.com/trek10inc/awsets/resource"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type AWSIoTThingGroup struct {
}

func init() {
	i := AWSIoTThingGroup{}
	listers = append(listers, i)
}

func (l AWSIoTThingGroup) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.IoTThingGroup}
}

func (l AWSIoTThingGroup) List(ctx context.AWSetsCtx) (*resource.Group, error) {

	svc := iot.New(ctx.AWSCfg)
	rg := resource.NewGroup()
	var nextToken *string
	for {
		groups, err := svc.ListThingGroupsRequest(&iot.ListThingGroupsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}).Send(ctx.Context)
		if err != nil {
			return rg, fmt.Errorf("failed to list iot thing groups: %w", err)
		}
		for _, group := range groups.ThingGroups {
			r := resource.New(ctx, resource.IoTThingGroup, group.GroupName, group.GroupName, group)
			rg.AddResource(r)
		}
		if groups.NextToken == nil {
			break
		}
		nextToken = groups.NextToken
	}
	return rg, nil
}
