package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redisと接続する
func RedisConnect() *redis.Client {
	//var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 1000,
	})
	return rdb
}

// Redisにデータをセットする
func RedisDataSet(rdb *redis.Client, key string, value string) {
	var ctx = context.Background()

	//err := rdb.Set(ctx, key, value, 0).Err()
	//err := rdb.Set(ctx, key, value, 30*time.Second).Err()
	err := rdb.Set(ctx, key, value, 10*time.Minute).Err()
	//err := rdb.Set(ctx, key, value, 1*time.Hour).Err()
	if err != nil {
		fmt.Println("登録エラー")
		fmt.Println(err)
	}
	fmt.Println("登録完了")
}

// Redisのデータを取得する
func RedisDataGet(rdb *redis.Client, key string) (ret string) {
	var ctx = context.Background()

	ret, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("取得エラー")
		fmt.Println(err)
		return
	}

	fmt.Println("Result: ", ret)

	// 登録されているKeyをすべて取得、表示する
	cmd := rdb.Keys(ctx, "*")
	res, err := cmd.Result()
	if err != nil {
		fmt.Println(err)
	} else {
		log.Println("Redisデータ登録数:", len(res))
		for i := 0; i < len(res); i++ {
			fmt.Println(res[i])
		}
	}

	return ret
}
