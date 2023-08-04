package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mertess/cryptocurrency-statistics-api-go/internal/worker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetLastPrice(ctx *gin.Context) {
	const postgresDsn = "host=localhost user=postgres password=56321 port=5432 dbname=CryptocurrencyStatistics sslmode=disable"
	db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		ctx.AbortWithError(500, err)
		return
	}

	currency, ok := ctx.GetQuery("currency")
	if !ok {
		ctx.AbortWithStatusJSON(400, "currency is required")
		return
	}

	var deal models.Deal
	db.Where("currencies = ?", currency).Order("updated_at_utc desc").First(&deal)

	if deal.Id == 0 {
		ctx.AbortWithStatus(404)
		return
	}

	ctx.JSON(200, deal)
}
