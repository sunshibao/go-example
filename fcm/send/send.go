package main

import (
	"context"
	"fmt"
	"go-example/fcm"

	"firebase.google.com/go/messaging"
)

func main() {
	topicOnly := &messaging.Message{
		//Topic: "2680007500000016408_en",
		Topic: "31274_zh_cn",

		Data: map[string]string{
			"packageName":   "com.miHoYo.GenshinImpact",
			"clickJumpType": "3",
		},
		Notification: &messaging.Notification{
			Title: "你预约的`包名`已安装", //你预约的包名已安装
			Body:  "你预约的`包名`已安装",
		},
	}
	ctx := context.Background()

	name, err := fcm.RyFcmClient.Send(ctx, topicOnly)
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err)
	}

}
