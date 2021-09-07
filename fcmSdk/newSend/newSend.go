package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func main() {
	var serviceAccountKey = []byte(`{
          "type": "service_account",
		  "project_id": "ry-push",
		  "private_key_id": "48d6c4e35f59ebd6f3f6506306cbe59c1f0ff749",
		  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDlkyXGCf2p+Vn5\n3DZkeSNv+8etdpNmS4zisC0sG2wn46sr0UMYqa/7HhBrgqSaaAeuaRykI1kPPSmu\nGg5vXhHlD7MZReLkYPTuhtiXhHmOjk6bLtnDGV24+p9mAud3czsCmhZRn0dwWH1C\nzxqA8Ns85UnhMSNqHWlKec4sZAUOT81DeKblqkMwlwUpqC50NaitdBECESGRkl3V\nqYzyVkECy5LhauI2hnuSt08nBism6WClRxIuOfjI0oPE4z25b/plBndx6KWhMqXV\nC+EI0FYzLWypHHjm127e4ElbVIeoSl8YYqj94bVfh8Jj7KP4jRDoFFAE43zvbxR8\nZnlJB1GjAgMBAAECggEABxixLQ8nq1vnsfINL7Y7krLD6DLx+SeyX1rWd/ZnCvqe\nKBjBYV9JbyftllfaUmLkgznWO/qoaPCukPusNjlrBdt3i/tIOWAH4jcCi03PMMlp\nrkwDSOUdVmJPkBMx6CmEiZ7ivBDXiqZUROmVplELrd+2lMKfQiCjnbzofcsftDf9\n9iO5d9ib9Bu5RF9DmkA/nC1K6h95uyh07hxNGFzgp1u6nzqvkrdQCIfHt3a3HV3b\n4Ua3GzhfF5MPnr1cG7Jz5zmpRxIqehNRjF9kgAtzrY/HyHpf+yRFVHgRcmehUPEL\n66rTBt093WRJBg6qlHdSj/+QUHKt9MVP8OAG+lPFxQKBgQD3kVQz45MtjQrwtKuv\nd7mjGL5ybN52IIW8J66FkqLxvfIhmYmtXcjvyIH+FUFLkoD5g0LMtTMAbZQcOSvL\n7CMiOabVGN5JFyJMZG6BGnToyrtLIaQT0yJvIess4XBYwnJrZhL8LaefMal5igiJ\nxn6yYkxG3umehE8P9W0rhjHyfwKBgQDtZO3jC5mgbSUAKbAxteKApsVyLNpmMj3s\nxpJu+YO1g0TQwcO+36m3iG+1oJghF06dHfCYCorhMiJDWHS4dJUSbGx7L8R5XN8C\nGHsPlxXQDrAXg3LcYxjq1giE8Y1HBmaN5r50OAvUAkChZELBZ8u6Schq3wrbOTJr\naLZ5SsYG3QKBgDZqTiy9l8sKVl5EB/ygf0A3Hx09isRCL8MEodqgOYqTKpZyDg36\nMEFsPA5iE6ENyaWOUW93YId2jniJpHPFKo/KRj5OogVEvXg3FwbvjsTgUryX312w\nKcBtnyiVQMFxs/6hSAj/6/kUzGB3k5rc81o4OvXU51q9UDd8xYssiuv3AoGBAJ2d\n6phZ5lluvidahoKq6cVTdTr2btd1uknQGf+WqQ1GJ9WXIRlFNVEHRGxKQVePOwH0\nk/7O2SDmAXvHak/iD+wYkvpDX+bYc7TXfjV+sdvfNKmX/BY5sZySGTvziULEDClh\nL8jIQYo1KKY/hFcXTFvWizvx9SpS0pggAS+NJuf9AoGAFsF7WXqd6wGyIDH43Gj0\n89Jm7eu8uMOuE0OeTsVV8jebBTPnUod1S6URn18IVTdv6SlGIY9ICMAZoZMf+cjj\nUXFldRMiDMRcDm+i+nV1u+/lRXL2yoA0WzO4hjkvm7soX26YR/mxg+s6PzF4YqY/\nh2+00xqVxk3lGY2NzZdXCAY=\n-----END PRIVATE KEY-----\n",
		  "client_email": "firebase-adminsdk-rgf08@ry-push.iam.gserviceaccount.com",
		  "client_id": "103297913384884974542",
		  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
		  "token_uri": "https://oauth2.googleapis.com/token",
		  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-rgf08%40ry-push.iam.gserviceaccount.com"
     }`)

	// opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	opt := option.WithCredentialsJSON(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	uid := "25696773511053390"
	// This registration token comes from the client FCM SDKs.
	registrationToken := "fxV6ToLJh3A:APA91bGfcaOl4mmnj_mPY7MTscjzT0aZLvyK5xaLLboWavxFoeqc3hZu_npEtaINebzHAfOrARg4kn9RmWC9ZYKvhqJrPhnNI43qtUruQsvd7Or7w_ZnDG4agOMM_7xB0J4ci9UHPT5S"

	// These registration tokens come from the client FCM SDKs.
	registrationTokens := []string{
		registrationToken,
		// ...
		"cBYSJNhfG_Q:APA91bFzLxiSVynUc2thc6aGfF1ba_6WoJvOctw2_1cIlUEr2r7Pf-n_Qk6uisLpc9Whcf-UU4WwcjnRwLTm_Zok1pH2RGw2_WvLmaT_AdZp84caH29haB4gQFIdrc0wQSr-vVgR0F3o",
	}

	subscribe(ctx, client, uid, registrationTokens)
	// sendMsgToToken(ctx, client, registrationToken)
	sendMsgToTopic(ctx, client, uid)
	// createCustomToken(ctx, app)
}

func sendMsgToTopic(ctx context.Context, client *messaging.Client, topic string) {
	notification := &messaging.Notification{
		Title: "0007 2019-04-23",
		Body:  "fffffffffffffffffffffffffff.",
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "88888",
			"time":  "2:45",
		},
		Notification: notification,
		Topic:        topic,
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}

func sendMsgToToken(ctx context.Context, client *messaging.Client, registrationToken string) {
	// See documentation on defining a message payload.

	notification := &messaging.Notification{
		Title: "$GOOG up 1.43% on the day",
		Body:  "$GOOG gained 11.80 points to close at 835.67, up 1.43% on the day.",
	}

	// timestampMillis := int64(12345)

	message := &messaging.Message{
		// Data: map[string]string{
		//  "score": "850",
		//  "time": "2:45",
		// },
		Notification: notification,
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: "title",
				Body:  "body",
				//      Icon: "icon",
			},
			FcmOptions: &messaging.WebpushFcmOptions{
				Link: "https://fcm.googleapis.com/",
			},
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}

func subscribe(ctx context.Context, client *messaging.Client, topic string, registrationTokens []string) {
	// Subscribe the devices corresponding to the registration tokens to the
	// topic.
	response, err := client.SubscribeToTopic(ctx, registrationTokens, topic)
	if err != nil {
		log.Fatalln(err)
	}
	// See the TopicManagementResponse reference documentation
	// for the contents of response.
	fmt.Println(response.SuccessCount, "tokens were subscribed successfully")
}

func createCustomToken(ctx context.Context, app *firebase.App) {
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := authClient.CustomToken(ctx, "25696773511053390")
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)
}
