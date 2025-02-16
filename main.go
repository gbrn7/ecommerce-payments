package main

import (
	"ecommerce-payments/cmd"
	"ecommerce-payments/helpers"
)

func main() {
	// load config
	helpers.SetupConfig()

	// load log
	helpers.SetupLogger()

	// load db
	helpers.SetupPostgreSQL()

	// run redis
	// helpers.SetupRedis()

	// run kafka consumer
	go cmd.ServeKafkaConsumerGroup()

	// run http
	cmd.ServeHTTP()
}
