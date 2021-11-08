package main

import (
	"context"
	"fmt"
	fcmTest "go-example/fcmSdk"
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
		Version: "test13-version",
	}
)

func main() {
	ctx := context.Background()
	client, err := fcmTest.NewClient(ctx, testMessagingConfig)
	if err != nil {
		fmt.Printf(err.Error())
	}
	client.FcmEndpoint = fcmTest.DefaultMessagingEndpoint
	topicOnly := &fcmTest.Message{
		Token: "test13-topic",
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
