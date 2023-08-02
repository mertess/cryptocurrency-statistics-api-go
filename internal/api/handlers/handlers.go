package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mertess/cryptocurrency-statistics-api-go/internal/worker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetLastPrice(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		rw.WriteHeader(400)
		rw.Write([]byte("Unsupported method"))
		return
	}

	const postgresDsn = "host=localhost user=postgres password=56321 port=5432 dbname=CryptocurrencyStatistics sslmode=disable"
	db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	currency := req.URL.Query().Get("currency")
	var deal models.Deal
	db.Where("currencies = ?", currency).Order("updated_at_utc desc").First(&deal)

	if deal.Id == 0 {
		rw.WriteHeader(404)
		return
	}

	response, err := json.Marshal(deal)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.Write(response)
}
