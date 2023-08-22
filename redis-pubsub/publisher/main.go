package main

import (
	"context"
	"log"
	"net/http"

	Customer "github.com/0xSherlokMo/messaging-patterns/redis-pubsub/customer"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("cannot ping redis, connection failed")
	}

	r := gin.Default()
	r.POST("/customer", func(c *gin.Context) {
		var customer Customer.Customer
		err := c.Bind(&customer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed on serialization",
			})
			return
		}

		serializedCustomer, err := customer.Json()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "serialization error",
			})
			return
		}

		_, err = rdb.Publish(ctx, Customer.ChannelCreation, serializedCustomer).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "cannot publish",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "processing",
		})
	})

	r.Run()
}
