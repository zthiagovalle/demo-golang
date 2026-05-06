package producers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/zthiagovalle/demo-golang/src/domain/models"
)

type ProductSNSProducer struct {
	client   *sns.Client
	topicARN string
}

func NewProductSNSProducer(client *sns.Client, topicARN string) *ProductSNSProducer {
	return &ProductSNSProducer{client: client, topicARN: topicARN}
}

type productCreatedEvent struct {
	EventType string         `json:"event_type"`
	Product   models.Product `json:"product"`
}

func (p *ProductSNSProducer) PublishCreated(ctx context.Context, product *models.Product) error {
	payload, err := json.Marshal(productCreatedEvent{
		EventType: "product.created",
		Product:   *product,
	})
	if err != nil {
		return err
	}

	_, err = p.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(p.topicARN),
		Message:  aws.String(string(payload)),
	})
	return err
}
