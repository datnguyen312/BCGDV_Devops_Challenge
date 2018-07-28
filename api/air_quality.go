package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CityAirQuality struct {
	City            string `json:"city"`
	AirQualityIndex int    `json:"air_quality_index"`
}

type AirVisualResponse struct {
	Data struct {
		City    string `json:"city"`
		Current struct {
			Pollution struct {
				Ts     time.Time `json:"ts"`
				Aqius  int       `json:"aqius"`
				Mainus string    `json:"mainus"`
				Aqicn  int       `json:"aqicn"`
				Maincn string    `json:"maincn"`
			} `json:"pollution"`
		} `json:"current"`
	} `json:"data"`
}

var apiKey = "XN6S5LmLBiGTZFHTq"

type AirQualityQueryer interface {
	FindAirQualityIndex(lat float64, lon float64) (CityAirQuality, error)
}

type AirVisualQueryer struct {
	client *http.Client
}

func (q *AirVisualQueryer) FindAirQualityIndex(lat float64, lon float64) (quality CityAirQuality, err error) {

	url := fmt.Sprintf("http://api.airvisual.com/v2/nearest_city?lat=%v&lon=%v&key=%v", lat, lon, apiKey)

	var res *http.Response
	res, err = q.client.Get(url)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	airVisualResponse := &AirVisualResponse{}
	err = json.Unmarshal(body, airVisualResponse)
	if err != nil {
		log.Fatal(err)
		return
	}

	quality = CityAirQuality{
		City:            airVisualResponse.Data.City,
		AirQualityIndex: airVisualResponse.Data.Current.Pollution.Aqius,
	}
	return
}
