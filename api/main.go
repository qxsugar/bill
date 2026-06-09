package main

import (
	"log"

	"github.com/qxsugar/bill/api/bootstrap"
	"github.com/qxsugar/bill/api/router"
)

func main() {
	bootstrap.InitDB()

	r := router.New()
	log.Println("starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
