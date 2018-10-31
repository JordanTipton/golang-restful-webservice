package main

import (
	"os"

	"github.com/jordantipton/golang-restful-webservice/app"
)

func main() {
	a := app.App{}
	a.Initialize(os.Getenv("DSN"))
	a.Run(":8080")
}
