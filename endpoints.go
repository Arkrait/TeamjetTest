package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type SortRequest struct {
	Array []int `json:"array"`
	Uniq bool `json:"uniq"`
}

type SortResponse struct {
	Array []int `json:"array"`
}

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64     `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code  int    `json:"cod"`
}

func SetupRoutes(r *chi.Mux) *chi.Mux {
	r.Route("/api", func(r chi.Router) {
		r.Get("/now", handleNow)
		r.Post("/sort", handleSort)
		r.Get("/weather", handleWeather)
	})

	return r
}

func handleNow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc).Format("2006-01-02 15:04:05 -0700 UTC")
	_, err = w.Write([]byte(now))
	if err != nil {
		panic(err)
	}
}

func handleSort(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var sortRequest SortRequest
	err := decoder.Decode(&sortRequest)
	if err != nil {
		panic(err)
	}
	arrayLength := len(sortRequest.Array)
	if arrayLength > 100 || arrayLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sort.Ints(sortRequest.Array)
	var resultArray = sortRequest.Array
	if sortRequest.Uniq {
		resultArray = sliceUniq(resultArray)
	}
	sortResponse := SortResponse{Array: resultArray}
	jsonData, err := json.Marshal(sortResponse)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func handleWeather(w http.ResponseWriter, req *http.Request) {
	cityName := req.URL.Query().Get("city")
	if cityName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := http.DefaultClient.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=e6aeb4dd03ba8a696a6378afd7e69a4d", cityName))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	var weatherResponse WeatherResponse
	err = decoder.Decode(&weatherResponse)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(strconv.Itoa(int(weatherResponse.Main.Temp))))
	if err != nil {
		panic(err)
	}
}

func sliceUniq(s []int) []int {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				// delete
				s = append(s[:j], s[j+1:]...)
				j--
			}
		}
	}
	return s
}
