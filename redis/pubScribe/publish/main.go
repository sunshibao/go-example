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
	n, err := client.Publish("chat", "hello").Result()
	if err != nil{
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("%d clients received the message\n", n)
}