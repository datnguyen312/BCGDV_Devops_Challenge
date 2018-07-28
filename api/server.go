package main

import (
	"time"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	storer, retriever, err := InitDatabase()
	if err != nil {
		log.Fatal("InitDatabase: ", err)
		os.Exit(1)
	}
	qualityqueryer := &AirVisualQueryer{client: &http.Client{}}
	StartServer(9090, CreateHandlers(storer, retriever, qualityqueryer))
}

func StartServer(port int, handler http.Handler) {
	fmt.Println("Server starting on port", port)
	s := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler, IdleTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second}
	log.Fatal(s.ListenAndServe())
}

func CreateHandlers(storer CityQueryStorer, retriever CityQueryRetriever, qualityqueryer AirQualityQueryer) http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/", handlePingRequest)
	m.Handle("/air-quality", createAirQualityQueryHandler(qualityqueryer, storer))
	m.Handle("/queried-cities", createQueriedCitiesHandler(retriever))
	return m
}

func handlePingRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Air Quality API running")
}

func createQueriedCitiesHandler(cqr CityQueryRetriever) http.Handler {

	var handler http.HandlerFunc

	handler = func(w http.ResponseWriter, r *http.Request) {
		queriedCities := cqr.RetrieveAllQueriedCities()
		responseJSON, err := json.Marshal(queriedCities)

		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}

	return handler
}

func createAirQualityQueryHandler(q AirQualityQueryer, s CityQueryStorer) http.Handler {

	var handler http.HandlerFunc

	handler = func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		latStr := query.Get("lat")
		if latStr == "" {
			latStr = "52"
		}
		lonStr := query.Get("lon")
		if lonStr == "" {
			lonStr = "0"
		}

		print(latStr, lonStr)

		lat, _ := strconv.ParseFloat(latStr, 32)
		lon, _ := strconv.ParseFloat(lonStr, 32)

		airQuality, err := q.FindAirQualityIndex(lat, lon)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.StoreCityQuery(airQuality.City)

		responseJSON, err := json.Marshal(airQuality)

		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}

	return handler

}
