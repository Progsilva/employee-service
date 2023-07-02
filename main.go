package main

import (
	"log"

	"github.com/Progsilva/employee-service/cmd/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	if err = a.Start(); err != nil {
		log.Fatal(err)
	}
}
