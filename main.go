package main

import (
	"fmt"

	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {

	database.ConnectDB()
	broker.ConnectRabbitMQ()
	fmt.Println("Hello World")
}
