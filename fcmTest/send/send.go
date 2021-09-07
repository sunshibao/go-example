package main

import (
	"context"
	"fmt"
	"go-example/fcmSdk"
	"go-example/fcmSdk/internal"

	"google.golang.org/api/option"
)

var (
	token               = "ya29.a0ARrdaM9Byh3pjmNbXLuQR7FFO_VHhleMcJpx12yQv0RP2gNGcMm0LK1BqJIS3Aqx95MRFkPma8EVKF4ak1D7e7rksGrmrxzqhvasoWW3DqB3iKMxaQSXodCwEu51xvfUtW8_nlRpcMS25K5TNiSG9j0MoP1l"
	testMessagingConfig = &internal.MessagingConfig{
		ProjectID: "ry-push",
		Opts: []option.ClientOption{
			option.WithTokenSource(&internal.MockTokenSource{AccessToken: token}),
		},
		Version: "test-version",
	}
)

func main() {
	ctx := context.Background()
	client, err := fcmSdk.NewClient(ctx, testMessagingConfig)
	if err != nil {
		fmt.Printf(err.Error())
	}
	client.FcmEndpoint = fcmSdk.DefaultMessagingEndpoint
	topicOnly := &fcmSdk.Message{
		Topic: "test-topic",
		Data: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}
	name, err := client.Send(ctx, topicOnly)
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err)
	}

}
