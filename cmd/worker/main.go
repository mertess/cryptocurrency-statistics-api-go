package main

import "github.com/mertess/cryptocurrency-statistics-api-go/internal/worker"

func main() {
	const dsn = "host=localhost user=postgres password=56321 port=5432 dbname=CryptocurrencyStatistics sslmode=disable"

	worker.Run(dsn)
}
