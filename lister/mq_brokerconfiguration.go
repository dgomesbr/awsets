package lister

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/trek10inc/awsets/context"
	"github.com/trek10inc/awsets/resource"
)

type AWSAmazonMQBrokerConfiguration struct {
}

func init() {
	i := AWSAmazonMQBrokerConfiguration{}
	listers = append(listers, i)
}

func (l AWSAmazonMQBrokerConfiguration) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.AmazonMQBrokerConfiguration}
}

func (l AWSAmazonMQBrokerConfiguration) List(ctx context.AWSetsCtx) (*resource.Group, error) {
	svc := mq.New(ctx.AWSCfg)
	rg := resource.NewGroup()

	var nextToken *string
	for {
		configurations, err := svc.ListConfigurationsRequest(&mq.ListConfigurationsInput{
			MaxResults: aws.Int64(100),
			NextToken:  nextToken,
		}).Send(ctx.Context)
		if err != nil {
			return rg, fmt.Errorf("failed to list mq broker configurations: %w", err)
		}
		for _, v := range configurations.Configurations {
			r := resource.New(ctx, resource.AmazonMQBrokerConfiguration, v.Id, v.Name, v)

			rg.AddResource(r)
		}
		if configurations.NextToken == nil {
			break
		}
		nextToken = configurations.NextToken
	}
	return rg, nil
}
