package publisher

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"

	"web-app/utils"
)

// see https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/pubsub/pubsub.go

var (
    Topic *pubsub.Topic
)

type PushRequest struct {
	Message struct {
        Attributes map[string]string
        Data       []byte
        ID         string `json:"message_id"`
    }
	Subscription string
}

func Init() {
    ctx := context.Background()

	client, err := pubsub.NewClient(ctx, utils.MustGetenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	topicName := utils.MustGetenv("PUBSUB_TOPIC")
	Topic = client.Topic(topicName)

	// Create the topic if it doesn't exist.
	exists, err := Topic.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
        log.Fatal(fmt.Sprintf("Topic %v doesn't exist", topicName))
	}

    log.Printf("[PUBSUB] up and ready.")
}

func CreatePushRequest() *PushRequest {
    return &PushRequest{}
}

func SendMessage(payload string) (string, error) {
	ctx := context.Background()

	msg := &pubsub.Message{
		Data: []byte(payload),
	}

	if _, err := Topic.Publish(ctx, msg).Get(ctx); err != nil {
		return "", err
	}

	log.Printf("[PUBSUB] <- %s", msg.Data)

	return "ok", nil
}
