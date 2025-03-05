package main

import (
	"github.com/maxpawgdbs/yandex-go/calculator"
	"github.com/maxpawgdbs/yandex-go/handlers"
	"log"
	"net/http"
)

func main() {
	calculator.Initial()
	log.Println("Starting Server")
	http.HandleFunc("/api/v1/calculate", handlers.CalculatorHandler)
	http.ListenAndServe(":8080", nil)
}
