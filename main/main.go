package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxpawgdbs/yandex-go/calculator"
	"github.com/maxpawgdbs/yandex-go/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	calculator.Initial()
	err := os.MkdirAll("database", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	log.Println("Starting Server")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", handlers.CalculatorHandler)
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionAnswer)
	http.ListenAndServe(":8080", r)
}
