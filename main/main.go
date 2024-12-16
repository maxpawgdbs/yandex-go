package main

import (
	"encoding/json"
	"fmt"
	"github.com/maxpawgdbs/yandex-go/calculator"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Expression string `json:"expression"`
}

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var req Request
		err = json.Unmarshal(body, &req)
		result, err := calculator.Calc(req.Expression)
		fmt.Fprint(w, result, err)
		log.Println("POST", string(body), req, result, err)
	} else {
		fmt.Fprint(w, "Only GET")
		w.WriteHeader(404)
		log.Println("GET 404")
	}
}

func main() {
	http.HandleFunc("/", CalculatorHandler)
	http.ListenAndServe(":8080", nil)
}
