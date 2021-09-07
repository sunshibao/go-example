package main

import (
	"context"
	"fmt"
	"go-example/fcm"

	"firebase.google.com/go/messaging"
)

func main() {
	ctx := context.Background()
	topicOnly := &messaging.Message{
		Topic: "test-topic",
		Data: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}
	name, err := fcm.FcmClient.Send(ctx, topicOnly)
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err)
	}

}
