package main

import (
	"fmt"
	log "github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

// keyspace

type PSubscribeCallback func (pattern, channel, message string)

type PSubscriber struct {
	client redis.PubSubConn
	cbMap map[string]PSubscribeCallback
}


func (c *PSubscriber) PConnect(ip string, port uint16) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	//conn, err := redis.Dial("tcp", ip + ":" + strconv.Itoa(int(port)))
	if err != nil {
		log.Critical("redis dial failed.")
	}

	c.client = redis.PubSubConn{conn}
	c.cbMap = make(map[string]PSubscribeCallback)

	go func() {
		for {
			log.Debug("wait...")
			switch res := c.client.Receive().(type) {
			case redis.Message:
				pattern := res.Pattern
				channel := string(res.Channel)
				message := string(res.Data)
				if (pattern == "__keyspace@0__:blog"){
					switch  message {
					case "set":
						// do something
						fmt.Println("set", channel)
					case "del":
						fmt.Println("del", channel)
					case "expire":
						fmt.Println("expire", channel)
					case "expired":
						fmt.Println("expired", channel)
					}
				}
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				log.Error("error handle...")
				continue
			}
		}
	}()

}
func (c *PSubscriber)Psubscribe(channel interface{}, cb PSubscribeCallback) {
	err := c.client.PSubscribe(channel)
	if err != nil{
		log.Critical("redis Subscribe error.")
	}

	c.cbMap[channel.(string)] = cb
}

func TestPubCallback(patter , chann, msg string){
	log.Debug( "TestPubCallback patter : " + patter + " channel : ", chann, " message : ", msg)
}

func main() {

	log.Info("===========main start============")


	var psub PSubscriber
	psub.PConnect("127.0.0.1", 6379)
	psub.Psubscribe("__keyspace@0__:blog", TestPubCallback)
	for{
		time.Sleep(1 * time.Second)
	}
}