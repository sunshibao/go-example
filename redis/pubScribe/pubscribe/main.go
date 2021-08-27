
package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pubsub := client.Subscribe("chat")
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}