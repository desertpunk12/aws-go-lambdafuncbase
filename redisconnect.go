package red

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"main/constant"
	"time"
)

var rd *redis.Client

func RD() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), constant.QueryTimeout*time.Second)
	defer cancel()

	if rd == nil || rd.Ping(ctx).Err() != nil {
		rd = redis.NewClient(&redis.Options{
			Addr:     constant.RDHOST + ":" + constant.RDPORT,
			Password: "",
			DB:       0,
		})

		log.Println("created new redis client!")
		if rd.Ping(ctx).Err() != nil {
			log.Println("redis connection error!", rd.Ping(ctx).Err())
		}
	} else {
		log.Println("reused redis client")
	}

	return rd
}
