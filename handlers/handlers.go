package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/maxpawgdbs/yandex-go/calculator"
	"github.com/maxpawgdbs/yandex-go/structs"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var req structs.Request
		err = json.Unmarshal(body, &req)
		//result, err := calculator.Calc(req.Expression)

		id := rand.Int()

		jsonResult, _ := json.Marshal(structs.ResponseResult{id, "proccessing", 0})
		os.WriteFile(fmt.Sprintf("database/%d.json", id), jsonResult, 0644)

		go calculator.Calc(req.Expression, id)

		jsonOut, _ := json.Marshal(structs.ResponseOK{Id: id})
		fmt.Fprint(w, string(jsonOut))
		log.Println("POST", req, string(jsonOut), 201)

	} else {
		w.WriteHeader(405)
		log.Println(r.Method, 405)
	}
}
func ExpressionAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		w.WriteHeader(404)
		jsonOut, _ := json.Marshal(map[string]structs.ResponseResult{"expression": structs.ResponseResult{id, "value error", 404}})
		fmt.Fprint(w, string(jsonOut))
		log.Println(string(jsonOut))
		return
	}
	_, err = os.Stat(fmt.Sprintf("database/%d.json", id))
	if os.IsNotExist(err) {
		w.WriteHeader(404)
		jsonOut, _ := json.Marshal(map[string]structs.ResponseResult{"expression": structs.ResponseResult{id, "not found", 404}})
		fmt.Fprint(w, string(jsonOut))
		log.Println(string(jsonOut))
		return
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("database/%d.json", id))
	if err != nil {
		w.WriteHeader(500)
		jsonOut, _ := json.Marshal(map[string]structs.ResponseResult{"expression": structs.ResponseResult{id, "some errors", 500}})
		fmt.Fprint(w, string(jsonOut))
		log.Println(string(jsonOut))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("{\"expression\": %s", string(data)))
	log.Println(string(data))
}
func ExpressionsList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("database")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, " чёто с бд сорян не будет кина")
		return
	}
	out := make([]structs.ResponseResult, 0)
	for _, file := range files {
		var structura structs.ResponseResult
		data, _ := ioutil.ReadFile(fmt.Sprintf("database/%s", file.Name()))
		json.Unmarshal(data, &structura)
		out = append(out, structura)
	}
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(map[string][]structs.ResponseResult{"expressions": out})
	fmt.Fprint(w, string(result))
	log.Println(string(result))
}

var OrkestratorGoroutinesCount int = 0
var COMPUTING_POWER int = 1000
var mu sync.Mutex

func Initial() {
	godotenv.Load(".env")
	value := os.Getenv("COMPUTING_POWER")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable COMPUTING_POWER")
			os.Exit(0)
		}
		COMPUTING_POWER = intvalue
	} else {
		COMPUTING_POWER = 1000
	}
}
func OrkestratorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		for {
			if OrkestratorGoroutinesCount < COMPUTING_POWER {
				mu.Lock()
				OrkestratorGoroutinesCount++
				mu.Unlock()
				body, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				var req structs.AgentResponse
				json.Unmarshal(body, &req)
				timer := time.NewTimer(time.Duration(req.Operation_time) * time.Millisecond)
				result := 0.0
				if req.Operation == "+" {
					result = req.Arg1 + req.Arg2
				} else if req.Operation == "-" {
					result = req.Arg1 - req.Arg2
				} else if req.Operation == "*" {
					result = req.Arg1 * req.Arg2
				} else if req.Operation == "/" {
					result = req.Arg1 / req.Arg2
				}
				<-timer.C
				w.WriteHeader(http.StatusOK)
				out, _ := json.Marshal(structs.AgentResult{result})
				fmt.Fprint(w, string(out))
				mu.Lock()
				OrkestratorGoroutinesCount--
				mu.Unlock()
				break
			}
		}
	}
}
