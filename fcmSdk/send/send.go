package main

import (
	"context"
	"fmt"
	"go-example/fcmSdk"
	"go-example/fcmSdk/internal"
	"time"

	"google.golang.org/api/option"
)

var (
	token               = "ya29.a0ARrdaM8UGnS_8kNBx0Lh1rwQ9X0MotbqRQuWfHYZozGlGzI4T5eaZAuk9p_V3gk1YOVApeBS0vQ5aJIjyoEXvHIxRVdnDIdS5yupxD8YRnQhnA-0PJa7yRnDiqG7xje22Ocj4VOo-9kNT1ubjKYYdmeXz32D"
	testMessagingConfig = &internal.MessagingConfig{
		ProjectID: "ry-push",
		Opts: []option.ClientOption{
			option.WithTokenSource(&internal.MockTokenSource{AccessToken: token}),
		},
		Version: "test-version",
	}

	ttlWithNanos = time.Duration(1500) * time.Millisecond
	ttl          = time.Duration(10) * time.Second
	invalidTTL   = time.Duration(-10) * time.Second

	badge           = 42
	badgeZero       = 0
	timestampMillis = int64(12345)
	timestamp       = time.Unix(0, 1546304523123*1000000).UTC()
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

	status, err := client.Send(ctx, topicOnly)
	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

}
