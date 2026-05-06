package consumers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/infra/gateways"
)

type ProductConsumer struct {
	sqsClient *sqs.Client
	queueURL  string
	catalog   *gateways.CatalogGateway
}

func NewProductConsumer(sqsClient *sqs.Client, queueURL string, catalog *gateways.CatalogGateway) *ProductConsumer {
	return &ProductConsumer{
		sqsClient: sqsClient,
		queueURL:  queueURL,
		catalog:   catalog,
	}
}

type productCreatedEvent struct {
	EventType string         `json:"event_type"`
	Product   models.Product `json:"product"`
}

func (c *ProductConsumer) Run(ctx context.Context) {
	log.Printf("[consumer] starting on queue %s", c.queueURL)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[consumer] stopping")
			return
		default:
		}

		out, err := c.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
		})
		if err != nil {
			log.Printf("[consumer] receive error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, msg := range out.Messages {
			c.handle(ctx, *msg.Body)
			_, _ = c.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(c.queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
}

func (c *ProductConsumer) handle(ctx context.Context, body string) {
	var event productCreatedEvent
	if err := json.Unmarshal([]byte(body), &event); err != nil {
		log.Printf("[consumer] invalid payload: %v", err)
		return
	}

	info, err := c.catalog.GetProductInfo(ctx, event.Product.ID)
	if err != nil {
		log.Printf("[consumer] catalog gateway error for product %s: %v", event.Product.ID, err)
		return
	}
	log.Printf("[consumer] enriched product %s: category=%s popularity=%d", event.Product.ID, info.Category, info.Popularity)
}
