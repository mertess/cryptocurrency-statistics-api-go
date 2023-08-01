package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/mertess/cryptocurrency-statistics-api-go/internal/worker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(postgresDsn string) {
	db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to postgres...")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println("Connected to redis...")

	buffer := make([]byte, 4096)
	for {
		resp, err := http.Get("https://yobit.net/api/3/ticker/eth_usd-btc_usd-eth_btc-trx_usdt")
		if err != nil {
			log.Println(err)
		}

		n, err := resp.Body.Read(buffer)
		if err != nil {
			log.Println(err)
		}
		buffer = buffer[:n]

		values := make(map[string]map[string]interface{})
		json.Unmarshal(buffer, &values)

		log.Println("Got deal:")
		log.Println(string(buffer))
		log.Println()

		for k, v := range values {
			cachedValue, err := redisClient.Get(k).Result()
			if err != nil && err != redis.Nil {
				log.Println(err)
				continue
			}

			value := v["last"].(float64)

			if cachedValue != "" {
				cachedValueFloat, err := strconv.ParseFloat(cachedValue, 64)
				if err != nil {
					log.Println(err)
					continue
				}

				if cachedValueFloat == value {
					continue
				} else {
					_, err = redisClient.Set(k, value, time.Minute*2).Result()
					if err != nil && err != redis.Nil {
						log.Println(err)
						continue
					}
				}
			} else {
				_, err = redisClient.Set(k, value, time.Minute*2).Result()
				if err != nil && err != redis.Nil {
					log.Println(err)
					continue
				}
			}

			db.Create(models.NewDeal(k, value, time.Now().UTC()))
		}

		time.Sleep(time.Second)
	}
}
