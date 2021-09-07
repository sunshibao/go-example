package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
)

func main() {
	ctx := context.Background()
	//config := &firebase.Config{ProjectID: "ry-push"}
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomToken(context.Background(), "some-uid")
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)
}
