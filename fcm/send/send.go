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
		Topic: "com.tencent.tmgp.sgame",
		Data: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
		Notification: &messaging.Notification{
			Title: "Title",
			Body:  "Body",
		},
	}
	name, err := fcm.RyFcmClient.Send(ctx, topicOnly)
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err)
	}

}
