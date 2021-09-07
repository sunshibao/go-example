package main

import (
	"context"
	"fmt"
	"go-example/fcm"

	"firebase.google.com/go/messaging"
)

func main() {
	ctx := context.Background()
	var testMessages = []*messaging.Message{
		{
			Topic: "topic1",
			Data: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
		{
			Topic: "topic2",
			Data: map[string]string{
				"k3": "v3",
				"k4": "v4",
			},
		},
	}
	br, err := fcm.RyFcmClient.SendAll(ctx, testMessages)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Printf("%+v\n", br)
	for _, v := range br.Responses {
		fmt.Printf("%+v\n", v)
	}
}
