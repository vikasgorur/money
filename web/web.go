package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vikasgorur/money"
)

type RateCache struct {
	UsdToInr  float64
	UpdatedAt time.Time
}

var cachedRate RateCache

// expired returns true if the cache is more than an hour old.
func expired() bool {
	return cachedRate.UpdatedAt.Before(time.Now().Add(-time.Duration(60) * time.Minute))
}

// getCachedRate returns the cached USD to INR conversion rate.
// It calls money.GetUsdToInr() to populate its cache.
func getCachedRate() float64 {
	if cachedRate.UsdToInr == 0 || expired() {
		log.Println("updating cached rate")
		cachedRate.UsdToInr = money.GetUsdToInr()
		cachedRate.UpdatedAt = time.Now()
	}

	return cachedRate.UsdToInr
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type ConvertResponse struct {
	Text string `json:"text"`
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("text")
	amount, err := money.ParseAmount(input)
	if err != nil {
		log.Println("invalid input: " + input)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output := amount.Convert(getCachedRate()).FormatValue()
	response, err := json.Marshal(ConvertResponse{Text: output})

	if err != nil {
		log.Println("json serialization failed for " + output)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	if err != nil {
		log.Println("couldn't write response: " + err.Error())
		return
	}

	log.Printf("converted '%s' to '%s'", input, output)
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/convert", handleConvert).Methods("GET")
	http.Handle("/", r)
}
