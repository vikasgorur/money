package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
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

func renderTemplate(w http.ResponseWriter, name string, value interface{}) {
	cwd, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(path.Join(cwd, "static", name))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, value)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

type ConvertResponse struct {
	Text string `json:"text"`
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("text")
	amount, err := money.Parse(input)
	if err != nil {
		log.Println("invalid input: " + input)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output := fmt.Sprint(amount.Convert(getCachedRate()))
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
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("static/css/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("static/js/"))))

	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/convert", handleConvert).Methods("GET")
	http.Handle("/", r)
}
