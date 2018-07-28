package main

import (
	"encoding/json"
	"github.com/dnaeon/go-vcr/recorder"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StorerDoubleFn func(city string) error

func (sd StorerDoubleFn) StoreCityQuery(city string) error {
	return sd(city)
}

type RetrieverDoubleFn func() []UserQuery

func (rd RetrieverDoubleFn) RetrieveAllQueriedCities() []UserQuery {
	return rd()
}

type AirQualityQueryDoubleFn func(lat float64, lon float64) (CityAirQuality, error)

func (fn AirQualityQueryDoubleFn) FindAirQualityIndex(lat float64, lon float64) (CityAirQuality, error) {
	return fn(lat, lon)
}

func Test_AirQualityQuery(t *testing.T) {

	var sd StorerDoubleFn
	sd = func(city string) error {
		return nil
	}
	var cqr RetrieverDoubleFn
	cqr = func() []UserQuery {
		return []UserQuery{}
	}

	r, err := recorder.New("fixtures/air-visual")
	if err != nil {
		t.Fatal(err)
		return
	}

	defer r.Stop()

	client := &http.Client{
		Transport: r,
	}

	q := &AirVisualQueryer{client: client}

	handlers := CreateHandlers(sd, cqr, q)

	server := httptest.NewServer(handlers)
	defer server.Close()

	var res *http.Response
	res, err = http.Get(server.URL + "/air-quality?lat=53.5&lon=-2.25")

	if err != nil {
		t.Fatal(err)
	}

	var airq CityAirQuality

	err = json.NewDecoder(res.Body).Decode(&airq)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Status code was incorrect, got: %d, want: %d.", res.StatusCode, 200)
	}

	if airq.City != "Manchester" {
		t.Errorf("City was incorrect, got: %d, want: %d.", airq.City, "Manchester")
	}

	if airq.AirQualityIndex != 19 {
		t.Errorf("AirQualityIndex was incorrect, got: %d, want: %d.", airq.AirQualityIndex, 19)
	}

}
