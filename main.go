package main

import (
	"fmt"

	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
)

func main() {

	broker.ConnectRabbitMQ()
	fmt.Println("Hello World")
}
