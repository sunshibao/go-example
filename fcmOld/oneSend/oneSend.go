package main

import (
	"fmt"
	fcm "myGo/fcmOld"
)

const (
	serverKey = "AAAAbCW6n1Q:APA91bFTctykFW5QkzDBXnt3nES7jgqD_JQ0VY1qUf8nffd3_nZrQplP3sqKf19HehL8Czs6zfI3RhKqhmlYqllRS0rQdbkgHr4o74--y7ZoXodQcMohEfiMxGJhaqu5EKPQa2AQZusf"
	topic     = "/topics/someTopic"
)

func main() {

	data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}

	c := fcm.NewFcmClient(serverKey)
	c.NewFcmMsgTo(topic, data)

	status, err := c.Send()

	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

}
