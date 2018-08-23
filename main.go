package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jordantipton/golang-restful-webservice/app"
)

func main() {
	// Load configuration
	if err := app.LoadConfig("./config/config.json"); err != nil {
		panic(fmt.Errorf("Failed to load configuration: %s", err))
	}

}
