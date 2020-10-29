package pubsub

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"cloud.google.com/go/pubsub"
)

func PullMsgsConcurrenyControl(w io.Writer, projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	// Must set ReceiveSettings.Synchronous to false (or leave as default) to enable
	// concurrency settings. Otherwise, NumGoroutines will be set to 1.
	sub.ReceiveSettings.Synchronous = false
	// NumGoroutines is the number of goroutines sub.Receive will spawn to pull messages concurrently.
	sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()

	// Receive messages for 10 seconds.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			messageData := string(msg.Data)
			fmt.Fprintf(w, "Got message :%q\n", messageData)
			log.Printf("[PUBSUB] -> %s", messageData)
			msg.Ack()
		}
	}()

	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}
