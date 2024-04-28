package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	Location struct {
		Error    bool   `json:"erro"`
		Cep      string `json:"cep"`
		Location string `json:"localidade"`
	}

	WeatherApiResp struct {
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
		Current Current `json:"current"`
	}

	Current struct {
		Temp_C float64 `json:"temp_c"`
		Temp_F float64 `json:"temp_f"`
	}

	ResponseDto struct {
		Temp_C float64 `json:"temp_C"`
		Temp_K float64 `json:"temp_K"`
		Temp_F float64 `json:"temp_F"`
	}
)

const (
	weatherApiKey = "58550850a19a4f7fb6a132328242104"
)

func main() {
	http.HandleFunc("GET /{cep}", Handle)
	http.ListenAndServe(":8080", nil)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	if valid := validCep(cep); !valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	respCep, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		http.Error(w, "can not get location", http.StatusInternalServerError)
		return
	}
	defer respCep.Body.Close()

	var l Location
	err = json.NewDecoder(respCep.Body).Decode(&l)
	if err != nil {
		http.Error(w, "can not decode location", http.StatusInternalServerError)
		return
	}

	if l.Error {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	openWeatherUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherApiKey, url.QueryEscape(l.Location))

	respWeather, err := http.Get(openWeatherUrl)
	if respWeather.StatusCode != http.StatusOK || err != nil {
		fmt.Println(openWeatherUrl)

		http.Error(w, "can not get weather", respWeather.StatusCode)
		return
	}
	defer respWeather.Body.Close()

	var weather WeatherApiResp
	if err = json.NewDecoder(respWeather.Body).Decode(&weather); err != nil {
		http.Error(w, "can not decode weather", http.StatusInternalServerError)
		return
	}
	currentTemp := getCurrentTemp(weather)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentTemp)
}

func validCep(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func getCurrentTemp(w WeatherApiResp) ResponseDto {
	return ResponseDto{
		Temp_C: w.Current.Temp_C,
		Temp_F: w.Current.Temp_F,
		Temp_K: w.Current.Temp_C + 273,
	}
}
