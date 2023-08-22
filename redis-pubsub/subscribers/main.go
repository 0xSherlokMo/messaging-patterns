package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	Customer "github.com/0xSherlokMo/messaging-patterns/redis-pubsub/customer"
	"github.com/redis/go-redis/v9"
)

var dbReplica = "MongoDB"

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	replica := os.Getenv("replica")
	if replica != "" {
		dbReplica = replica
	}
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("cannot ping redis, connection failed")
	}

	closer := make(chan struct{})
	go ProcessCreatingCustomer(ctx, rdb)

	log.Printf("%s replica server is up \n", dbReplica)
	<-closer
}

func ProcessCreatingCustomer(ctx context.Context, redisClient *redis.Client) {
	pubsub := redisClient.Subscribe(ctx, Customer.ChannelCreation)
	for message := range pubsub.Channel() {
		log.Printf("[%s] [NEW MESSAGE] %s \n", dbReplica, message.String())
		var customer Customer.Customer
		err := json.Unmarshal([]byte(message.Payload), &customer)
		if err != nil {
			log.Printf("[%s] [ERROR couldn't serialize message %s \n", dbReplica, message.String())
		}

		// do some processing to replicate data
		log.Printf("[%s] [ACK] message %s processed successfully, customer: %+v", dbReplica, message.String(), customer)
	}
}
