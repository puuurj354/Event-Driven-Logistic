package main

import (
	"fmt"

	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	_ = database.InitPostgres(cfg.Database.URL)

	fmt.Println("Order Service is running...")

}
