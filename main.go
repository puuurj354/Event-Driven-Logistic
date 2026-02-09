package main

import (
	"fmt"

	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {

	database.ConnectDB()
	fmt.Println("Hello World")
}
