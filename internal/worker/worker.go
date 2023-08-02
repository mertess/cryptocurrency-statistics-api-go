package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mertess/cryptocurrency-statistics-api-go/internal/cache"
	"github.com/mertess/cryptocurrency-statistics-api-go/internal/worker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	redisHost = "localhost:6379"
	password  = ""
	dbNumber  = 0
)

type DealsJson map[string]map[string]interface{}

func Run(postgresDsn string) {
	db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to postgres...")

	cache := cache.NewRedis(redisHost, password, dbNumber, time.Hour*48)
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

		values := make(DealsJson)
		json.Unmarshal(buffer, &values)

		log.Println("Got deal:")
		log.Println(string(buffer))
		log.Println()

		createDeals(values, cache, db)

		time.Sleep(time.Minute)
	}
}

func createDeals(values DealsJson, cache *cache.RedisCache, db *gorm.DB) {
	for k, v := range values {
		floatValue := v["last"].(float64)
		value, ok := cache.GetFloat64(k)
		if ok && value == floatValue {
			continue
		}

		ok = cache.SetFloat64(k, floatValue)
		if ok {
			db.Create(models.NewDeal(k, floatValue, time.Now().UTC()))
		}
	}
}
